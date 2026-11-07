"use client";

import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { deleteFile, downloadFile, FileItem, formatFileSize } from "@/lib/api";

interface FileTableProps {
  files: FileItem[];
  loading: boolean;
  onRefresh: () => void;
  onUnauthorized: () => void;
}

export default function FileTable({
  files,
  loading,
  onRefresh,
  onUnauthorized,
}: FileTableProps) {
  const [search, setSearch] = useState("");
  const [busyId, setBusyId] = useState<string | null>(null);
  const [errMsg, setErrMsg] = useState<string | null>(null);

  const filtered = files.filter((f) =>
    f.name.toLowerCase().includes(search.toLowerCase())
  );

  const onDownload = async (file: FileItem) => {
    setErrMsg(null);
    setBusyId(file.id);
    try {
      await downloadFile(file.id, file.name);
    } catch (e) {
      const m = e instanceof Error ? e.message : "Gagal mengunduh file";
      if (m === "UNAUTHORIZED") return onUnauthorized();
      setErrMsg(m);
    } finally {
      setBusyId(null);
    }
  };

  const onDelete = async (file: FileItem) => {
    setErrMsg(null);
    setBusyId(file.id);
    try {
      await deleteFile(file.id);
      onRefresh();
    } catch (e) {
      const m = e instanceof Error ? e.message : "Gagal menghapus file";
      if (m === "UNAUTHORIZED") return onUnauthorized();
      setErrMsg(m);
    } finally {
      setBusyId(null);
    }
  };

  const fmtDate = (iso: string) => {
    const d = new Date(iso);
    return isNaN(d.getTime())
      ? iso || "-"
      : d.toLocaleString("id-ID", {
          day: "2-digit",
          month: "short",
          year: "numeric",
          hour: "2-digit",
          minute: "2-digit",
        });
  };

  return (
    <div className="rounded-2xl border border-slate-800 bg-slate-900/70 p-6 shadow-xl">
      <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 className="text-lg font-semibold text-white">Daftar File Anda</h2>
          <p className="text-sm text-slate-400">
            Total {files.length} file aktif tersimpan
          </p>
        </div>

        <div className="flex items-center gap-3">
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Cari nama file..."
            className="w-full rounded-xl border border-slate-700 bg-slate-950 px-3.5 py-2 text-sm text-slate-100 placeholder-slate-500 focus:border-indigo-500 focus:outline-none sm:w-64"
          />
          <button
            type="button"
            onClick={onRefresh}
            disabled={loading}
            className="rounded-xl border border-slate-700 bg-slate-800 px-3.5 py-2 text-sm font-medium text-slate-200 hover:bg-slate-700 disabled:opacity-50"
          >
            Refresh
          </button>
        </div>
      </div>

      {errMsg && (
        <div className="mb-4 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
          {errMsg}
        </div>
      )}

      {loading ? (
        <div className="py-14 text-center text-sm text-slate-400">
          Memuat daftar file...
        </div>
      ) : filtered.length === 0 ? (
        <div className="rounded-xl border border-slate-800 bg-slate-950/40 py-14 text-center text-sm text-slate-400">
          {search ? "File tidak ditemukan." : "Belum ada file yang diunggah."}
        </div>
      ) : (
        <div className="overflow-x-auto rounded-xl border border-slate-800">
          <table className="w-full border-collapse text-left text-sm">
            <thead>
              <tr className="border-b border-slate-800 bg-slate-950/70 text-xs uppercase text-slate-400">
                <th className="px-4 py-3.5">Nama File</th>
                <th className="px-4 py-3.5">Tipe MIME</th>
                <th className="px-4 py-3.5">Ukuran</th>
                <th className="px-4 py-3.5">Waktu Upload</th>
                <th className="px-4 py-3.5 text-right">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/70">
              <AnimatePresence>
                {filtered.map((f) => (
                  <motion.tr
                    key={f.id}
                    initial={{ opacity: 0, y: 8 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, x: -16, transition: { duration: 0.2 } }}
                    className="hover:bg-slate-800/40"
                  >
                    <td className="px-4 py-3.5 font-medium text-slate-100">
                      {f.name}
                    </td>
                    <td className="px-4 py-3.5 font-mono text-xs text-slate-300">
                      {f.mime_type || "application/octet-stream"}
                    </td>
                    <td className="px-4 py-3.5 text-slate-300">
                      {formatFileSize(f.size)}
                    </td>
                    <td className="px-4 py-3.5 text-slate-400">
                      {fmtDate(f.created_at)}
                    </td>
                    <td className="px-4 py-3.5 text-right">
                      <div className="inline-flex gap-2">
                        <motion.button
                          whileHover={{ scale: 1.05 }}
                          whileTap={{ scale: 0.95 }}
                          type="button"
                          disabled={busyId === f.id}
                          onClick={() => onDownload(f)}
                          className="rounded-lg border border-indigo-500/30 bg-indigo-500/10 px-3 py-1.5 text-xs font-medium text-indigo-300 hover:bg-indigo-500/20 disabled:opacity-50"
                        >
                          Download
                        </motion.button>
                        <motion.button
                          whileHover={{ scale: 1.05 }}
                          whileTap={{ scale: 0.95 }}
                          type="button"
                          disabled={busyId === f.id}
                          onClick={() => onDelete(f)}
                          className="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-1.5 text-xs font-medium text-red-300 hover:bg-red-500/20 disabled:opacity-50"
                        >
                          Hapus
                        </motion.button>
                      </div>
                    </td>
                  </motion.tr>
                ))}
              </AnimatePresence>
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
