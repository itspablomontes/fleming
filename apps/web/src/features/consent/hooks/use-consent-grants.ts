import { useQuery } from "@tanstack/react-query";

import { getActiveGrants, getMyGrants } from "../api";

export const useMyConsentGrants = () =>
	useQuery({
		queryKey: ["consent", "grants", "grantor"],
		queryFn: getMyGrants,
	});

export const useActiveConsentGrants = () =>
	useQuery({
		queryKey: ["consent", "grants", "active"],
		queryFn: getActiveGrants,
	});
