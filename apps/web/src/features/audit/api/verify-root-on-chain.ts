import { apiClient } from "@/lib/api-client";

export interface VerifyRootOnChainResponse {
	anchored: boolean;
	timestamp: number;
	blockNumber: number | null;
	txHash: string | null;
}

export const verifyRootOnChain = async (
	root: string,
): Promise<VerifyRootOnChainResponse> => {
	return apiClient(`/api/audit/verify/${encodeURIComponent(root)}`);
};

