import axios from "axios";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8000";

export async function login(email: string, password: string) {
  const res = await axios.post(
    `${API_BASE}/auth/login/`,
    { email, password },
    { withCredentials: true }
  );
  return res.data;
}

export async function signup(
  username: string,
  email: string,
  password: string,
  password2: string
) {
  const res = await axios.post(
    `${API_BASE}/auth/registration/`,
    { username, email, password, password2 },
    { withCredentials: true }
  );
  return res.data;
}

export async function logout() {
  await axios.post(`${API_BASE}/auth/logout/`, {}, { withCredentials: true });
}

export async function googleLogin(accessToken: string) {
  const res = await axios.post(
    `${API_BASE}/auth/google/`,
    { access_token: accessToken },
    { withCredentials: true }
  );
  return res.data;
}
