import { apiClient } from "@/lib/api-client";

export interface FileKeyResponse {
  wrappedKey: string;
}

export interface GetFileKeyPayload {
  eventId: string;
  fileId: string;
  patientId?: string;
}

export const getFileKey = (payload: GetFileKeyPayload): Promise<FileKeyResponse> => {
  const query = payload.patientId ? `?patientId=${payload.patientId}` : "";
  return apiClient(
    `/api/timeline/events/${payload.eventId}/files/${payload.fileId}/key${query}`,
  );
};

export interface ShareFilePayload {
  eventId: string;
  fileId: string;
  grantee: string;
  wrappedKey: string;
}

export const shareFileKey = (payload: ShareFilePayload) => {
  return apiClient(`/api/timeline/events/${payload.eventId}/files/${payload.fileId}/share`, {
    method: "POST",
    body: {
      grantee: payload.grantee,
      wrappedKey: payload.wrappedKey,
    },
  });
};
