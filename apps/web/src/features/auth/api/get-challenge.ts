import { apiClient } from "@/lib/api-client";

export interface ChallengeRequest {
	address: string;
	domain: string;
	uri: string;
	chainId: number;
}

export interface ChallengeResponse {
	message: string;
}

export const getChallenge = async (data: ChallengeRequest): Promise<string> => {
	const res = await apiClient("/api/auth/challenge", { body: data });
	return res.message;
};
