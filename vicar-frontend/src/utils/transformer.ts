export function transformSize(size: number): string {
  if (size < 1024) {
    return size + " B";
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(2) + " KB";
  } else if (size < 1024 * 1024 * 1024) {
    return (size / (1024 * 1024)).toFixed(2) + " MB";
  } else {
    return (size / (1024 * 1024 * 1024)).toFixed(2) + " GB";
  }
}

export function transformFileType(filename: string): 'image' | 'video' | 'audio' | 'file' {
  const parts = filename.split(".");
  if (parts.length < 2) {
    return "file";
  }
  const ext = parts[parts.length - 1].toLowerCase();
  switch (ext) {
    case "png":
    case "jpg":
    case "jpeg":
    case "gif":
    case "bmp":
    case "webp":
      return "image";
    case "mp4":
    case "avi":
    case "mkv":
    case "mov":
    case "wmv":
    case "flv":
    case "mpeg":
    case "mpg":
    case "3gp":
    case "m4v":
    case "webm":
      return "video";
    case "mp3":
    case "wav":
    case "ogg":
    case "flac":
    case "aac":
      return "audio";
    default:
      return "file";
  }
}

export function transformFileTypeIcon(filename: string): string {
  const type = transformFileType(filename);
  switch (type) {
    case "image":
      return "fas fa-image";
    case "video":
      return "fas fa-video";
    case "audio":
      return "fas fa-music";
    default:
      return "fas fa-file";
  }
}