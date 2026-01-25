import { apiClient } from "@/lib/api-client";

interface MerkleRootResponse {
	root: string;
}

export const getMerkleRoot = async (batchId: string): Promise<MerkleRootResponse> => {
	return apiClient(`/api/audit/merkle/${batchId}`);
};
