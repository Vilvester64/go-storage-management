package utils

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/h2non/filetype"
)

var (
	ErrUnsupportedFileType = errors.New("tipe atau isi file tidak diizinkan (terdeteksi tidak valid atau berbahaya)")
	ErrEmptyFile           = errors.New("file yang diunggah kosong")
)

// daftar ekstensi berbahaya yang dilarang diunggah.
var blockedExtensions = map[string]bool{
	".sh":   true,
	".bash": true,
	".zsh":  true,
	".php":  true,
	".phtml": true,
	".exe":  true,
	".dll":  true,
	".so":   true,
	".dylib": true,
	".bat":  true,
	".cmd":  true,
	".ps1":  true,
	".vbs":  true,
	".js":   true,
	".mjs":  true,
	".py":   true,
	".rb":   true,
	".pl":   true,
	".cgi":  true,
	".jar":  true,
	".msi":  true,
	".com":  true,
	".scr":  true,
	".bin":  true,
	".elf":  true,
}

// daftar MIME hasil deteksi magic bytes yang dilarang (executable/script).
var blockedMIMEs = map[string]bool{
	"application/x-executable":      true,
	"application/x-msdownload":      true,
	"application/x-elf":             true,
	"application/x-mach-binary":     true,
	"application/x-sh":              true,
	"application/x-shellscript":     true,
	"application/vnd.microsoft.portable-executable": true,
}

// daftar ekstensi teks murni yang diizinkan jika magic bytes tidak terdeteksi sebagai binary.
var allowedTextExtensions = map[string]string{
	".txt":  "text/plain",
	".csv":  "text/csv",
	".md":   "text/markdown",
	".json": "application/json",
	".xml":  "application/xml",
	".yaml": "application/x-yaml",
	".yml":  "application/x-yaml",
	".log":  "text/plain",
}

// SanitizeFileName membersihkan nama file asli dari komponen path untuk mencegah Path Traversal.
func SanitizeFileName(filename string) string {
	clean := filepath.Base(strings.TrimSpace(filename))
	clean = strings.ReplaceAll(clean, "\x00", "")
	if clean == "." || clean == "/" || clean == "\\" || clean == "" {
		return "unnamed_file"
	}
	return clean
}

// GenerateSystemFileName menghasilkan nama unik berbasis UUID beserta ekstensi asli yang aman.
func GenerateSystemFileName(originalFilename string) string {
	ext := strings.ToLower(filepath.Ext(SanitizeFileName(originalFilename)))
	return uuid.NewString() + ext
}

// containsDangerousScriptSignature memeriksa apakah header berisi signature script (misal shebang #! atau tag PHP)
// meskipun file telah di-rename menjadi .txt atau ekstensi lain.
func containsDangerousScriptSignature(header []byte) bool {
	trimmed := bytes.TrimLeft(header, " \t\r\n\xef\xbb\xbf")
	lower := bytes.ToLower(trimmed)

	if bytes.HasPrefix(trimmed, []byte("#!")) {
		return true
	}
	if bytes.HasPrefix(lower, []byte("<?php")) || bytes.HasPrefix(lower, []byte("<?=")) {
		return true
	}
	if bytes.HasPrefix(lower, []byte("<script")) {
		return true
	}
	// Deteksi ELF (\x7fELF) atau MZ (Windows PE executable)
	if bytes.HasPrefix(header, []byte{0x7F, 'E', 'L', 'F'}) || bytes.HasPrefix(header, []byte{'M', 'Z'}) {
		return true
	}
	return false
}

// DetectAndValidateMIME membaca magic bytes (261 byte pertama) menggunakan h2non/filetype
// dan memvalidasi apakah file aman untuk disimpan.
func DetectAndValidateMIME(r io.Reader, originalFilename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(SanitizeFileName(originalFilename)))
	if blockedExtensions[ext] {
		return "", ErrUnsupportedFileType
	}

	header := make([]byte, 512)
	n, err := r.Read(header)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	if n == 0 {
		return "", ErrEmptyFile
	}
	header = header[:n]

	// Cegah script berbahaya (seperti .sh atau .php yang di-rename menjadi .txt/.jpg)
	if containsDangerousScriptSignature(header) {
		return "", ErrUnsupportedFileType
	}

	kind, err := filetype.Match(header)
	if err != nil {
		return "", err
	}

	if kind != filetype.Unknown {
		mimeType := kind.MIME.Value
		if blockedMIMEs[mimeType] {
			return "", ErrUnsupportedFileType
		}
		return mimeType, nil
	}

	// Jika filetype.Unknown, periksa apakah merupakan file teks murni yang valid dan diizinkan.
	if textMIME, ok := allowedTextExtensions[ext]; ok {
		if bytes.IndexByte(header, 0x00) != -1 || !utf8.Valid(header) {
			return "", ErrUnsupportedFileType
		}
		return textMIME, nil
	}

	return "", ErrUnsupportedFileType
}
