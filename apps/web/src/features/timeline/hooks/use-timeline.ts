import { useQuery } from "@tanstack/react-query";
import { getTimeline } from "../api";

export function useTimeline() {
	return useQuery({
		queryKey: ["timeline"],
		queryFn: () => getTimeline(),
		staleTime: 1000 * 60 * 5, // 5 minutes
	});
}
