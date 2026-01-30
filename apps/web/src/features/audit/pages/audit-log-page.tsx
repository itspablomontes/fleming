import { Link } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { type JSX, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";

import { getAuditLogs, queryAuditEntries, verifyIntegrity } from "../api";
import {
	type AuditFilterState,
	AuditFilters,
} from "../components/audit-filters";
import { AuditLogViewer } from "../components/audit-log-viewer";
import { IntegrityStatus } from "../components/integrity-status";
import { MerkleProofViewer } from "../components/merkle-proof-viewer";

const defaultFilters: AuditFilterState = {
	actor: "",
	resourceId: "",
	resourceType: "",
	action: "",
	startTime: "",
	endTime: "",
};

const toIsoString = (value: string): string | undefined => {
	if (!value) {
		return undefined;
	}
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) {
		return undefined;
	}
	return date.toISOString();
};

export function AuditLogPage(): JSX.Element {
	const [filters, setFilters] = useState<AuditFilterState>(defaultFilters);
	const [integrityStatus, setIntegrityStatus] = useState<{
		valid?: boolean;
		message?: string;
	}>({});

	const queryParams = useMemo(() => {
		return {
			actor: filters.actor || undefined,
			resourceId: filters.resourceId || undefined,
			resourceType: filters.resourceType || undefined,
			action: filters.action || undefined,
			startTime: toIsoString(filters.startTime),
			endTime: toIsoString(filters.endTime),
		};
	}, [filters]);

	const hasFilters = useMemo(() => {
		return Object.values(queryParams).some((value) => value);
	}, [queryParams]);

	const auditQuery = useQuery({
		queryKey: ["audit-logs", queryParams],
		queryFn: () =>
			hasFilters ? queryAuditEntries(queryParams) : getAuditLogs(),
	});

	const verifyMutation = useMutation({
		mutationFn: verifyIntegrity,
		onSuccess: (data) => {
			setIntegrityStatus({ valid: data.valid, message: data.message });
		},
		onError: (error) => {
			setIntegrityStatus({
				valid: false,
				message: error instanceof Error ? error.message : "Verification failed",
			});
		},
	});

	return (
		<div className="min-h-screen bg-gray-50 px-4 py-6 dark:bg-gray-950 md:px-8 md:py-10">
			<div className="mx-auto flex w-full max-w-5xl flex-col gap-6">
				<div className="flex flex-col gap-3">
					<Button
						variant="ghost"
						size="icon"
						asChild
						className="h-9 w-9 shrink-0"
						aria-label="Back to timeline"
					>
						<Link to="/">
							<ArrowLeft className="h-4 w-4" />
						</Link>
					</Button>
					<div className="flex flex-col gap-2">
						<h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
							Audit Log
						</h1>
					<p className="text-sm text-muted-foreground">
						Review the cryptographic audit trail for your account activity.
					</p>
					</div>
				</div>

				<AuditFilters
					filters={filters}
					onChange={setFilters}
					onReset={() => setFilters(defaultFilters)}
				/>

				<div className="grid gap-4 md:grid-cols-2">
					<IntegrityStatus
						isValid={integrityStatus.valid}
						message={integrityStatus.message}
					/>
					<div className="flex flex-col gap-2 rounded-lg border border-gray-200 bg-white p-4 dark:border-gray-800 dark:bg-gray-900">
						<p className="text-sm font-medium text-gray-900 dark:text-gray-100">
							Verify audit chain
						</p>
						<p className="text-xs text-muted-foreground">
							Run a full hash-chain verification for the latest entries.
						</p>
						<Button
							onClick={() => verifyMutation.mutate()}
							disabled={verifyMutation.isPending}
							variant="outline"
						>
							{verifyMutation.isPending ? "Verifying..." : "Verify integrity"}
						</Button>
					</div>
				</div>

				<MerkleProofViewer />

				<AuditLogViewer
					entries={auditQuery.data ?? []}
					isLoading={auditQuery.isLoading}
					errorMessage={
						auditQuery.error ? "Failed to load audit logs." : undefined
					}
					onRefresh={() => auditQuery.refetch()}
				/>
			</div>
		</div>
	);
}
