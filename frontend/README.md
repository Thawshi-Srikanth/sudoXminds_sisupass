# SiSu Pass front end (pnpm)

This is a Next.js project bootstrapped with `create-next-app`.

## 1. Prerequisites

- Node.js 18+ ([Download Node.js](https://nodejs.org/))
- pnpm ([Install pnpm](https://pnpm.io/installation))
- Git
- Running Django backend (for API endpoints)

## 2. Move to directory

```bash
cd frontend
```

## 3. Install Dependencies

```bash
pnpm install
```

## 4. Set Up Environment Variables

Create a `.env.local` file:

```env
# Django API
NEXT_PUBLIC_API_URL=http://localhost:8000/api

# Google Auth
NEXT_PUBLIC_GOOGLE_CLIENT_ID=<your-google-client-id>
```

## 5. Run Development Server

```bash
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser. Hot-reloading is enabled.

## 6. Recommended Packages

- `axios` & `apollo` for API calls
- `tailwindcss` for styling

## 7. Google OAuth Setup

1. Go to Google Cloud Console.
2. Enable OAuth 2.0 Client IDs..
3. Add the client ID and secret to `.env.local`.
