import type { JSX } from "react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
	useApproveConsent,
	useDenyConsent,
	useRevokeConsent,
} from "../hooks/use-consent-actions";
import {
	useActiveConsentGrants,
	useMyConsentGrants,
} from "../hooks/use-consent-grants";
import type { ConsentGrant, ConsentState } from "../types";
import { ConsentGrantCard } from "./consent-grant-card";
import { ConsentRequestCard } from "./consent-request-card";

const consentStates: Array<{ value: ConsentState | "all"; label: string }> = [
	{ value: "all", label: "All states" },
	{ value: "requested", label: "Requested" },
	{ value: "approved", label: "Approved" },
	{ value: "denied", label: "Denied" },
	{ value: "revoked", label: "Revoked" },
	{ value: "expired", label: "Expired" },
];

interface ConsentFilters {
	state: ConsentState | "all";
	search: string;
	startDate: string;
	endDate: string;
}

const defaultFilters: ConsentFilters = {
	state: "all",
	search: "",
	startDate: "",
	endDate: "",
};

const matchesFilters = (
	grant: ConsentGrant,
	filters: ConsentFilters,
): boolean => {
	if (filters.state !== "all" && grant.state !== filters.state) {
		return false;
	}
	if (filters.search) {
		const query = filters.search.toLowerCase();
		if (
			!grant.grantor.toLowerCase().includes(query) &&
			!grant.grantee.toLowerCase().includes(query)
		) {
			return false;
		}
	}
	if (filters.startDate) {
		const start = new Date(filters.startDate);
		if (!Number.isNaN(start.getTime()) && grant.createdAt < start) {
			return false;
		}
	}
	if (filters.endDate) {
		const end = new Date(filters.endDate);
		if (!Number.isNaN(end.getTime()) && grant.createdAt > end) {
			return false;
		}
	}
	return true;
};

export function ConsentDashboard(): JSX.Element {
	const [filters, setFilters] = useState<ConsentFilters>(defaultFilters);

	const myGrantsQuery = useMyConsentGrants();
	const activeGrantsQuery = useActiveConsentGrants();

	const approveMutation = useApproveConsent({
		onSuccess: async () => {
			await Promise.all([myGrantsQuery.refetch(), activeGrantsQuery.refetch()]);
			toast.success("Access approved.");
		},
		onError: (error) => {
			const message =
				error instanceof Error ? error.message : "Failed to approve consent";
			toast.error(message);
		},
	});

	const denyMutation = useDenyConsent({
		onSuccess: async () => {
			await Promise.all([myGrantsQuery.refetch(), activeGrantsQuery.refetch()]);
			toast.success("Access denied.");
		},
		onError: (error) => {
			const message =
				error instanceof Error ? error.message : "Failed to deny consent";
			toast.error(message);
		},
	});

	const revokeMutation = useRevokeConsent({
		onSuccess: async () => {
			await Promise.all([myGrantsQuery.refetch(), activeGrantsQuery.refetch()]);
			toast.success("Access revoked.");
		},
		onError: (error) => {
			const message =
				error instanceof Error ? error.message : "Failed to revoke consent";
			toast.error(message);
		},
	});

	const filteredMyGrants = useMemo(() => {
		return (myGrantsQuery.data ?? []).filter((grant) =>
			matchesFilters(grant, filters),
		);
	}, [myGrantsQuery.data, filters]);

	const filteredActiveGrants = useMemo(() => {
		return (activeGrantsQuery.data ?? []).filter((grant) =>
			matchesFilters(grant, filters),
		);
	}, [activeGrantsQuery.data, filters]);

	const isExpiredGrant = (grant: ConsentGrant) =>
		Boolean(grant.expiresAt && grant.expiresAt.getTime() <= Date.now());

	const pendingRequests = filteredMyGrants.filter(
		(grant) => grant.state === "requested",
	);

	const activeIssued = filteredMyGrants.filter(
		(grant) => grant.state === "approved" && !isExpiredGrant(grant),
	);

	const historyIssued = filteredMyGrants.filter((grant) => {
		if (grant.state === "requested") {
			return false;
		}
		if (grant.state === "approved") {
			return isExpiredGrant(grant);
		}
		return true;
	});

	const isLoading = myGrantsQuery.isLoading || activeGrantsQuery.isLoading;
	const hasError = Boolean(myGrantsQuery.error || activeGrantsQuery.error);

	return (
		<div className="space-y-6">
			<Card className="border-border bg-white dark:bg-gray-900">
				<CardHeader className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
					<div className="space-y-2">
						<CardTitle className="text-base font-semibold text-foreground">
							Consent dashboard
						</CardTitle>
						<p className="text-sm text-muted-foreground">
							Track requests, active grants, and history in one place.
						</p>
					</div>
					<Button
						variant="outline"
						onClick={() => {
							myGrantsQuery.refetch();
							activeGrantsQuery.refetch();
						}}
						disabled={isLoading}
					>
						Refresh
					</Button>
				</CardHeader>
				<CardContent className="space-y-5">
					<div className="grid gap-4 md:grid-cols-4">
						<div className="space-y-2">
							<Label htmlFor="consent-state-filter">State</Label>
							<Select
								value={filters.state}
								onValueChange={(value) =>
									setFilters((prev) => ({
										...prev,
										state: value as ConsentState | "all",
									}))
								}
							>
								<SelectTrigger id="consent-state-filter" className="w-full">
									<SelectValue placeholder="All states" />
								</SelectTrigger>
								<SelectContent>
									{consentStates.map((option) => (
										<SelectItem key={option.value} value={option.value}>
											{option.label}
										</SelectItem>
									))}
								</SelectContent>
							</Select>
						</div>
						<div className="space-y-2">
							<Label htmlFor="consent-search">Address</Label>
							<Input
								id="consent-search"
								value={filters.search}
								onChange={(event) =>
									setFilters((prev) => ({
										...prev,
										search: event.target.value,
									}))
								}
								placeholder="0x..."
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="consent-start">From</Label>
							<Input
								id="consent-start"
								type="date"
								value={filters.startDate}
								onChange={(event) =>
									setFilters((prev) => ({
										...prev,
										startDate: event.target.value,
									}))
								}
							/>
						</div>
						<div className="space-y-2">
							<Label htmlFor="consent-end">To</Label>
							<Input
								id="consent-end"
								type="date"
								value={filters.endDate}
								onChange={(event) =>
									setFilters((prev) => ({
										...prev,
										endDate: event.target.value,
									}))
								}
							/>
						</div>
					</div>

					{hasError && (
						<p role="alert" className="text-sm text-destructive">
							Failed to load consent data. Please refresh.
						</p>
					)}

					<Tabs defaultValue="pending">
						<TabsList className="w-full justify-start">
							<TabsTrigger value="pending">Pending requests</TabsTrigger>
							<TabsTrigger value="active">Active grants</TabsTrigger>
							<TabsTrigger value="history">History</TabsTrigger>
						</TabsList>

						<TabsContent value="pending" className="space-y-4">
							{isLoading ? (
								<div className="space-y-3">
									<Skeleton className="h-28 w-full" />
									<Skeleton className="h-28 w-full" />
								</div>
							) : pendingRequests.length === 0 ? (
								<p className="text-sm text-muted-foreground">
									No pending requests right now.
								</p>
							) : (
								pendingRequests.map((grant) => (
									<ConsentRequestCard
										key={grant.id}
										grant={grant}
										onApprove={approveMutation.mutate}
										onDeny={denyMutation.mutate}
										isApproving={approveMutation.isPending}
										isDenying={denyMutation.isPending}
									/>
								))
							)}
						</TabsContent>

						<TabsContent value="active" className="space-y-6">
							<div className="space-y-3">
								<p className="text-xs uppercase tracking-wide text-muted-foreground">
									Grants you issued
								</p>
								{isLoading ? (
									<div className="space-y-3">
										<Skeleton className="h-28 w-full" />
									</div>
								) : activeIssued.length === 0 ? (
									<p className="text-sm text-muted-foreground">
										No active grants you issued.
									</p>
								) : (
									activeIssued.map((grant) => (
										<ConsentGrantCard
											key={grant.id}
											grant={grant}
											onRevoke={revokeMutation.mutate}
											isRevoking={revokeMutation.isPending}
										/>
									))
								)}
							</div>

							<div className="space-y-3">
								<p className="text-xs uppercase tracking-wide text-muted-foreground">
									Grants issued to you
								</p>
								{isLoading ? (
									<div className="space-y-3">
										<Skeleton className="h-28 w-full" />
									</div>
								) : filteredActiveGrants.length === 0 ? (
									<p className="text-sm text-muted-foreground">
										No active grants assigned to you.
									</p>
								) : (
									filteredActiveGrants.map((grant) => (
										<ConsentGrantCard key={grant.id} grant={grant} />
									))
								)}
							</div>
						</TabsContent>

						<TabsContent value="history" className="space-y-3">
							{isLoading ? (
								<div className="space-y-3">
									<Skeleton className="h-28 w-full" />
									<Skeleton className="h-28 w-full" />
								</div>
							) : historyIssued.length === 0 ? (
								<p className="text-sm text-muted-foreground">
									No historical grants to show.
								</p>
							) : (
								historyIssued.map((grant) => (
									<ConsentGrantCard key={grant.id} grant={grant} />
								))
							)}
						</TabsContent>
					</Tabs>
				</CardContent>
			</Card>
		</div>
	);
}
