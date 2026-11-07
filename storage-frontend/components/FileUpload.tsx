"use client";

import React, { useRef, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { formatFileSize, uploadFile } from "@/lib/api";

interface FileUploadProps {
  onUploadSuccess: () => void;
  onUnauthorized: () => void;
}

export default function FileUpload({
  onUploadSuccess,
  onUnauthorized,
}: FileUploadProps) {
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setError(null);
    setSuccessMsg(null);
    if (e.target.files?.[0]) setSelectedFile(e.target.files[0]);
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
    setError(null);
    setSuccessMsg(null);
    if (e.dataTransfer.files?.[0]) setSelectedFile(e.dataTransfer.files[0]);
  };

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedFile) {
      setError("Silakan pilih file terlebih dahulu.");
      return;
    }

    setUploading(true);
    setError(null);
    setSuccessMsg(null);

    try {
      const uploaded = await uploadFile(selectedFile);
      setSuccessMsg(`File "${uploaded.name}" berhasil diunggah!`);
      setSelectedFile(null);
      if (fileInputRef.current) fileInputRef.current.value = "";
      onUploadSuccess();
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Gagal mengunggah file";
      if (msg === "UNAUTHORIZED") return onUnauthorized();
      setError(msg);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="rounded-2xl border border-slate-800 bg-slate-900/70 p-6 shadow-xl backdrop-blur">
      <div className="mb-4">
        <h2 className="text-lg font-semibold text-white">Unggah File Baru</h2>
        <p className="text-sm text-slate-400">
          Format aman didukung (Gambar, PDF, Dokumen, Arsip, Teks). Maksimal 50MB.
        </p>
      </div>

      <form onSubmit={handleUpload} className="space-y-4">
        <motion.div
          whileHover={{ scale: 1.005 }}
          whileTap={{ scale: 0.995 }}
          onDragOver={(e) => {
            e.preventDefault();
            setIsDragging(true);
          }}
          onDragLeave={(e) => {
            e.preventDefault();
            setIsDragging(false);
          }}
          onDrop={handleDrop}
          onClick={() => fileInputRef.current?.click()}
          className={`flex cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed p-8 text-center transition-all ${
            isDragging
              ? "border-indigo-500 bg-indigo-500/10"
              : selectedFile
                ? "border-emerald-500/60 bg-emerald-500/5"
                : "border-slate-700 bg-slate-950/50 hover:border-indigo-500/60 hover:bg-slate-900"
          }`}
        >
          <input
            ref={fileInputRef}
            type="file"
            onChange={handleFileChange}
            className="hidden"
          />

          <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-indigo-500/10 text-indigo-400">
            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              />
            </svg>
          </div>

          {selectedFile ? (
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className="space-y-1"
            >
              <p className="text-sm font-medium text-emerald-400">{selectedFile.name}</p>
              <p className="text-xs text-slate-400">
                Ukuran: {formatFileSize(selectedFile.size)}
              </p>
              <p className="pt-1 text-xs text-indigo-400 underline">
                Klik untuk mengganti file
              </p>
            </motion.div>
          ) : (
            <div className="space-y-1">
              <p className="text-sm font-medium text-slate-200">
                Klik untuk memilih file atau tarik & lepas ke sini
              </p>
              <p className="text-xs text-slate-500">
                Validasi otomatis dengan deteksi Magic Bytes MIME
              </p>
            </div>
          )}
        </motion.div>

        <AnimatePresence>
          {error && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              className="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300"
            >
              {error}
            </motion.div>
          )}

          {successMsg && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              className="rounded-xl border border-emerald-500/30 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-300"
            >
              {successMsg}
            </motion.div>
          )}
        </AnimatePresence>

        <div className="flex items-center justify-end gap-3">
          {selectedFile && !uploading && (
            <button
              type="button"
              onClick={() => {
                setSelectedFile(null);
                if (fileInputRef.current) fileInputRef.current.value = "";
              }}
              className="rounded-xl border border-slate-700 px-4 py-2.5 text-sm font-medium text-slate-300 transition hover:bg-slate-800"
            >
              Batal
            </button>
          )}

          <motion.button
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            type="submit"
            disabled={!selectedFile || uploading}
            className="rounded-xl bg-indigo-600 px-5 py-2.5 text-sm font-semibold text-white shadow-lg shadow-indigo-600/25 transition hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {uploading ? "Mengunggah..." : "Upload File"}
          </motion.button>
        </div>
      </form>
    </div>
  );
}
