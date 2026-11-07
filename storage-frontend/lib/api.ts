export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

const TOKEN_KEY = "jwt_token";
const USER_EMAIL_KEY = "user_email";

export interface ApiResponse<T = unknown> {
  status: "success" | "error";
  message: string;
  data?: T;
}

export interface FileItem {
  id: string;
  name: string;
  size: number;
  mime_type?: string;
  created_at: string;
}

export interface UserProfile {
  id: string;
  email: string;
}

export function saveToken(token: string, email?: string): void {
  if (typeof window === "undefined") return;
  localStorage.setItem(TOKEN_KEY, token);
  if (email) localStorage.setItem(USER_EMAIL_KEY, email);
  document.cookie = `${TOKEN_KEY}=${encodeURIComponent(token)}; path=/; max-age=86400; SameSite=Lax`;
}

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  const localToken = localStorage.getItem(TOKEN_KEY);
  if (localToken) return localToken;

  const match = document.cookie
    .split("; ")
    .find((row) => row.startsWith(`${TOKEN_KEY}=`));
  return match ? decodeURIComponent(match.split("=")[1]) : null;
}

export function getSavedEmail(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(USER_EMAIL_KEY);
}

export function clearToken(): void {
  if (typeof window === "undefined") return;
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_EMAIL_KEY);
  document.cookie = `${TOKEN_KEY}=; path=/; max-age=0; SameSite=Lax`;
}

/**
 * Helper pembungkus fetch bawaan yang otomatis menggunakan NEXT_PUBLIC_API_URL
 * dan menyematkan header Authorization: Bearer <token> jika token ada di Cookies / Storage.
 */
export async function fetchApi(
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> {
  const cleanEndpoint = endpoint.startsWith("/") ? endpoint : `/${endpoint}`;
  const url = `${API_BASE_URL}${cleanEndpoint}`;

  const headers = new Headers(options.headers || {});
  const token = getToken();

  if (token && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  return fetch(url, {
    ...options,
    headers,
  });
}

export async function loginUser(email: string, password: string): Promise<string> {
  const response = await fetchApi("/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const result: ApiResponse<{ token: string }> & { token?: string } =
    await response.json();

  if (!response.ok || result.status === "error") {
    throw new Error(result.message || "Email atau password salah");
  }

  const token = result.data?.token || result.token;
  if (!token) throw new Error("Token JWT tidak ditemukan pada response server");

  saveToken(token, email);
  return token;
}

export async function registerUser(email: string, password: string): Promise<UserProfile> {
  const response = await fetchApi("/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const result: ApiResponse<UserProfile> = await response.json();
  if (!response.ok || result.status === "error" || !result.data) {
    throw new Error(result.message || "Gagal mendaftarkan akun");
  }
  return result.data;
}

export async function fetchFiles(): Promise<FileItem[]> {
  const response = await fetchApi("/files", { method: "GET" });

  if (response.status === 401) {
    clearToken();
    throw new Error("UNAUTHORIZED");
  }

  const result: ApiResponse<FileItem[]> = await response.json();
  if (!response.ok || result.status === "error") {
    throw new Error(result.message || "Gagal mengambil daftar file");
  }
  return Array.isArray(result.data) ? result.data : [];
}

export async function uploadFile(file: File): Promise<FileItem> {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetchApi("/files/upload", {
    method: "POST",
    body: formData,
  });

  if (response.status === 401) {
    clearToken();
    throw new Error("UNAUTHORIZED");
  }

  const result: ApiResponse<FileItem> = await response.json();
  if (!response.ok || result.status === "error" || !result.data) {
    throw new Error(result.message || "Gagal mengunggah file");
  }
  return result.data;
}

export async function downloadFile(fileId: string, fileName: string): Promise<void> {
  const response = await fetchApi(`/files/download/${fileId}`, { method: "GET" });

  if (response.status === 401) {
    clearToken();
    throw new Error("UNAUTHORIZED");
  }

  if (!response.ok) {
    let errorMessage = "Gagal mengunduh file";
    try {
      const errJson: ApiResponse = await response.json();
      if (errJson.message) errorMessage = errJson.message;
    } catch {
      // ignore non-JSON body
    }
    throw new Error(errorMessage);
  }

  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.setAttribute("download", fileName);
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
}

export async function deleteFile(fileId: string): Promise<void> {
  const response = await fetchApi(`/files/${fileId}`, { method: "DELETE" });

  if (response.status === 401) {
    clearToken();
    throw new Error("UNAUTHORIZED");
  }

  const result: ApiResponse = await response.json();
  if (!response.ok || result.status === "error") {
    throw new Error(result.message || "Gagal menghapus file");
  }
}

export function formatFileSize(bytes: number): string {
  if (!bytes || bytes <= 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  const size = bytes / Math.pow(1024, i);
  return `${size.toFixed(i === 0 ? 0 : 2)} ${units[i]}`;
}
