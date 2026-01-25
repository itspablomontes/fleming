import { apiClient } from "@/lib/api-client";

import type { MerkleProof } from "../types";

interface VerifyMerkleResponse {
	valid: boolean;
}

export interface VerifyMerkleParams {
	root: string;
	entryHash: string;
	proof: MerkleProof;
}

export const verifyMerkleProof = async (
	params: VerifyMerkleParams,
): Promise<VerifyMerkleResponse> => {
	return apiClient("/api/audit/merkle/verify", { body: params });
};
