import { apiClient } from "@/lib/api-client";

export interface MultipartStartPayload {
  eventId: string;
  fileName: string;
  mimeType: string;
}

export interface MultipartStartResponse {
  uploadId: string;
  objectName: string;
}

export const startMultipartUpload = (
  payload: MultipartStartPayload,
): Promise<MultipartStartResponse> => {
  return apiClient(`/api/timeline/events/${payload.eventId}/files/multipart/start`, {
    method: "POST",
    body: {
      fileName: payload.fileName,
      mimeType: payload.mimeType,
    },
  });
};

export interface MultipartPartPayload {
  eventId: string;
  uploadId: string;
  objectName: string;
  partNumber: number;
  part: Blob;
}

export interface MultipartPartResponse {
  etag: string;
}

export const uploadMultipartPart = async (
  payload: MultipartPartPayload,
): Promise<MultipartPartResponse> => {
  const formData = new FormData();
  formData.append("uploadId", payload.uploadId);
  formData.append("objectName", payload.objectName);
  formData.append("partNumber", payload.partNumber.toString());
  formData.append("part", payload.part);

  const response = await fetch(
    `/api/timeline/events/${payload.eventId}/files/multipart/part`,
    {
      method: "PUT",
      body: formData,
      credentials: "include",
    },
  );

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || "Failed to upload part");
  }

  return response.json();
};

export interface MultipartCompletePayload {
  eventId: string;
  uploadId: string;
  objectName: string;
  fileName: string;
  mimeType: string;
  fileSize: number;
  wrappedKey: string;
  chunkSize: number;
  totalSize: number;
  ivLength: number;
  parts: Array<{ partNumber: number; etag: string }>;
}

export const completeMultipartUpload = (
  payload: MultipartCompletePayload,
) => {
  return apiClient(
    `/api/timeline/events/${payload.eventId}/files/multipart/complete`,
    {
      method: "POST",
      body: {
        uploadId: payload.uploadId,
        objectName: payload.objectName,
        fileName: payload.fileName,
        mimeType: payload.mimeType,
        fileSize: payload.fileSize,
        wrappedKey: payload.wrappedKey,
        chunkSize: payload.chunkSize,
        totalSize: payload.totalSize,
        ivLength: payload.ivLength,
        parts: payload.parts,
      },
    },
  );
};
