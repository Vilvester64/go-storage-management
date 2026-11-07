"use client";

import React, { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
  clearToken,
  fetchFiles,
  FileItem,
  formatFileSize,
  getSavedEmail,
  getToken,
} from "@/lib/api";
import { motion } from "framer-motion";
import FileUpload from "@/components/FileUpload";
import FileTable from "@/components/FileTable";

export default function DashboardPage() {
  const router = useRouter();
  const [files, setFiles] = useState<FileItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [authChecked, setAuthChecked] = useState(false);
  const [userEmail, setUserEmail] = useState<string | null>(null);

  const handleUnauthorized = useCallback(() => {
    clearToken();
    router.replace("/login");
  }, [router]);

  const loadFiles = useCallback(async () => {
    setLoading(true);
    try {
      const data = await fetchFiles();
      setFiles(data);
    } catch (err) {
      if (err instanceof Error && err.message === "UNAUTHORIZED") {
        handleUnauthorized();
        return;
      }
    } finally {
      setLoading(false);
    }
  }, [handleUnauthorized]);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.replace("/login");
      return;
    }
    setUserEmail(getSavedEmail());
    setAuthChecked(true);
    loadFiles();
  }, [router, loadFiles]);

  const handleLogout = () => {
    clearToken();
    router.push("/login");
  };

  if (!authChecked) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-950 text-slate-200">
        <div className="flex flex-col items-center gap-3">
          <span className="h-8 w-8 animate-spin rounded-full border-2 border-indigo-500 border-t-transparent" />
          <p className="text-sm text-slate-400">Memverifikasi sesi...</p>
        </div>
      </div>
    );
  }

  const totalStorageUsed = files.reduce((acc, item) => acc + (item.size || 0), 0);

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      {/* Top Navbar */}
      <header className="sticky top-0 z-20 border-b border-slate-800/80 bg-slate-950/80 backdrop-blur-md">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-4 sm:px-6">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-600/20 border border-indigo-500/30 text-indigo-400">
              <svg
                className="h-5 w-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z"
                />
              </svg>
            </div>
            <div>
              <h1 className="text-base font-bold text-white sm:text-lg">
                Storage Manager Dashboard
              </h1>
              {userEmail && (
                <p className="text-xs text-slate-400">{userEmail}</p>
              )}
            </div>
          </div>

          <button
            type="button"
            onClick={handleLogout}
            className="inline-flex items-center gap-2 rounded-xl border border-slate-700 bg-slate-900 px-4 py-2 text-xs font-semibold text-slate-200 transition hover:border-red-500/40 hover:bg-red-500/10 hover:text-red-300 sm:text-sm"
          >
            <svg
              className="h-4 w-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
              />
            </svg>
            Logout
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="mx-auto max-w-6xl space-y-6 px-4 py-8 sm:px-6">
        {/* Ringkasan Statistik */}
        <motion.div
          initial={{ opacity: 0, y: 15 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4 }}
          className="grid grid-cols-1 gap-4 sm:grid-cols-3"
        >
          <motion.div
            whileHover={{ y: -3 }}
            transition={{ duration: 0.2 }}
            className="rounded-2xl border border-slate-800 bg-slate-900/60 p-5 shadow-lg backdrop-blur"
          >
            <p className="text-xs font-medium uppercase tracking-wider text-slate-400">
              Total File Aktif
            </p>
            <p className="mt-2 text-2xl font-bold text-white">{files.length}</p>
          </motion.div>

          <motion.div
            whileHover={{ y: -3 }}
            transition={{ duration: 0.2 }}
            className="rounded-2xl border border-slate-800 bg-slate-900/60 p-5 shadow-lg backdrop-blur"
          >
            <p className="text-xs font-medium uppercase tracking-wider text-slate-400">
              Total Penyimpanan Terpakai
            </p>
            <p className="mt-2 text-2xl font-bold text-indigo-400">
              {formatFileSize(totalStorageUsed)}
            </p>
          </motion.div>

          <motion.div
            whileHover={{ y: -3 }}
            transition={{ duration: 0.2 }}
            className="rounded-2xl border border-slate-800 bg-slate-900/60 p-5 shadow-lg backdrop-blur"
          >
            <p className="text-xs font-medium uppercase tracking-wider text-slate-400">
              Batas Ukuran per Upload
            </p>
            <p className="mt-2 text-2xl font-bold text-emerald-400">50 MB</p>
          </motion.div>
        </motion.div>

        {/* Komponen Upload File */}
        <motion.div
          initial={{ opacity: 0, y: 15 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.1, duration: 0.4 }}
        >
          <FileUpload
            onUploadSuccess={loadFiles}
            onUnauthorized={handleUnauthorized}
          />
        </motion.div>

        {/* Komponen Tabel Daftar File */}
        <motion.div
          initial={{ opacity: 0, y: 15 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2, duration: 0.4 }}
        >
          <FileTable
            files={files}
            loading={loading}
            onRefresh={loadFiles}
            onUnauthorized={handleUnauthorized}
          />
        </motion.div>
      </main>
    </div>
  );
}
