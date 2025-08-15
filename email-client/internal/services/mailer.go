package services

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"sync"
	"text/template"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/google/uuid"
	"sisupass.com/email-client/internal/types"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type MailerService struct {
	dialer      *mail.Dialer
	sender      string
	emailQueue  chan *EmailJob
	emailStatus map[string]*types.EmailStatus
	statusMutex sync.RWMutex
	workers     int
	stats       *EmailStats
	statsMutex  sync.RWMutex
}

type EmailStats struct {
	TotalSent    int64     `json:"total_sent"`
	TotalFailed  int64     `json:"total_failed"`
	TotalPending int64     `json:"total_pending"`
	TotalQueued  int64     `json:"total_queued"`
	QueueSize    int       `json:"queue_size"`
	Workers      int       `json:"workers"`
	StartTime    time.Time `json:"start_time"`
	LastActivity time.Time `json:"last_activity"`
}

type EmailJob struct {
	ID      string
	Request *types.EmailRequest
	Retry   int
}

func NewMailerService(host string, port int, username, password, sender string, workers int) *MailerService {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 10 * time.Second

	service := &MailerService{
		dialer:      dialer,
		sender:      sender,
		emailQueue:  make(chan *EmailJob, 1000),
		emailStatus: make(map[string]*types.EmailStatus),
		workers:     workers,
		stats: &EmailStats{
			Workers:   workers,
			StartTime: time.Now(),
		},
	}

	for i := 0; i < workers; i++ {
		go service.worker()
	}

	return service
}

func (m *MailerService) GetStats() *EmailStats {
	m.statsMutex.RLock()
	defer m.statsMutex.RUnlock()

	stats := *m.stats
	stats.QueueSize = len(m.emailQueue)
	return &stats
}

func (m *MailerService) updateStats(sent, failed, pending, queued int64) {
	m.statsMutex.Lock()
	defer m.statsMutex.Unlock()

	m.stats.TotalSent += sent
	m.stats.TotalFailed += failed
	m.stats.TotalPending += pending
	m.stats.TotalQueued += queued
	m.stats.LastActivity = time.Now()
}

func (m *MailerService) worker() {
	for job := range m.emailQueue {
		err := m.sendEmailNow(job)

		m.statusMutex.Lock()
		if status, exists := m.emailStatus[job.ID]; exists {
			if err != nil {
				status.Status = "failed"
				status.FailReason = err.Error()
				m.updateStats(0, 1, 0, 0)

				if job.Retry < 3 {
					job.Retry++
					go func() {
						time.Sleep(time.Duration(job.Retry) * time.Minute)
						m.emailQueue <- job
					}()
					status.Status = "retrying"
				}
			} else {
				status.Status = "sent"
				now := time.Now()
				status.SentAt = &now
				m.updateStats(1, 0, 0, 0)
			}
		}
		m.statusMutex.Unlock()
	}
}

func (m *MailerService) sendEmailNow(job *EmailJob) error {
	req := job.Request

	templateFile := fmt.Sprintf("templates/%s.tmpl", req.Template)
	tmpl, err := template.ParseFS(templateFS, templateFile)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	subject := new(bytes.Buffer)
	if req.Subject != "" {
		subject.WriteString(req.Subject)
	} else {
		err = tmpl.ExecuteTemplate(subject, "subject", req.Data)
		if err != nil {
			return fmt.Errorf("failed to execute subject template: %v", err)
		}
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", req.Data)
	if err != nil {
		return fmt.Errorf("failed to execute plain body template: %v", err)
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", req.Data)
	if err != nil {
		return fmt.Errorf("failed to execute html body template: %v", err)
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", req.To)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	switch req.Priority {
	case types.PriorityHigh:
		msg.SetHeader("X-Priority", "1")
	case types.PriorityLow:
		msg.SetHeader("X-Priority", "5")
	default:
		msg.SetHeader("X-Priority", "3")
	}

	return m.dialer.DialAndSend(msg)
}

func (m *MailerService) GetEmailStatus(id string) (*types.EmailStatus, error) {
	m.statusMutex.RLock()
	defer m.statusMutex.RUnlock()

	status, exists := m.emailStatus[id]
	if !exists {
		return nil, fmt.Errorf("email not found")
	}

	return status, nil
}

func (m *MailerService) SendBulkEmail(req *types.BulkEmailRequest) ([]*types.EmailResponse, error) {
	responses := make([]*types.EmailResponse, 0, len(req.Recipients))

	for _, recipient := range req.Recipients {
		emailReq := &types.EmailRequest{
			To:       recipient,
			Template: req.Template,
			Subject:  req.Subject,
			Data:     req.Data,
			Priority: types.PriorityNormal,
		}
		resp, err := m.SendEmail(emailReq)
		if err != nil {
			log.Printf("Failed to queue email for %s: %v", recipient, err)
			continue
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

func (m *MailerService) SendEmail(req *types.EmailRequest) (*types.EmailResponse, error) {
	id := uuid.New().String()

	status := &types.EmailStatus{
		ID:        id,
		Status:    "pending",
		Recipient: req.To,
		Template:  req.Template,
		CreatedAt: time.Now(),
	}

	m.statusMutex.Lock()
	m.emailStatus[id] = status
	m.statusMutex.Unlock()

	job := &EmailJob{
		ID:      id,
		Request: req,
		Retry:   0,
	}

	m.updateStats(0, 0, 1, 1)

	if req.ScheduledAt != nil && req.ScheduledAt.After(time.Now()) {
		go func() {
			time.Sleep(time.Until(*req.ScheduledAt))
			m.emailQueue <- job
		}()
	} else {
		select {
		case m.emailQueue <- job:
		default:
			return nil, fmt.Errorf("email queue is full")
		}
	}

	return &types.EmailResponse{
		ID:        id,
		Status:    "queued",
		Message:   "Email queued for sending",
		SentAt:    time.Now(),
		Recipient: req.To,
	}, nil
}
