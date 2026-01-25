import { apiClient } from "@/lib/api-client";

export interface MerkleBatch {
	id: string;
	rootHash: string;
	startTime: string;
	endTime: string;
	entryCount: number;
	createdAt: string;
}

interface BuildMerkleResponse {
	batch: MerkleBatch;
	root: string;
}

export interface BuildMerkleParams {
	startTime?: Date | string;
	endTime?: Date | string;
}

const toTimestamp = (value?: Date | string) => {
	if (!value) {
		return undefined;
	}
	return value instanceof Date ? value.toISOString() : value;
};

export const buildMerkleTree = async (
	params: BuildMerkleParams,
): Promise<BuildMerkleResponse> => {
	return apiClient("/api/audit/merkle/build", {
		body: {
			startTime: toTimestamp(params.startTime),
			endTime: toTimestamp(params.endTime),
		},
	});
};
