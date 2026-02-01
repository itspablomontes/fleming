import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import {
	ArrowLeft,
	ChevronDown,
	FileText,
	Filter,
	RefreshCw,
	Upload,
} from "lucide-react";
import { type JSX, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/ui/empty-state";
import { cn } from "@/lib/utils";

import { listMerkleBatches, queryAuditEntries } from "../api";
import {
	type AuditFilterState,
	AuditFilters,
} from "../components/audit-filters";
import { AuditLogTable } from "../components/audit-log-table";
import { AuditStatusBar } from "../components/audit-status-bar";
import { MerkleTreeViewer } from "../components/merkle-tree-viewer";
import { PaginationControls } from "../components/pagination-controls";
import { VerifyRootWidget } from "../components/verify-root-widget";
import { useAuditStore } from "../stores/audit-store";
import type { AuditAction, AuditTargetType } from "../types";

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

// Quick filter chips for common action types
// Map "Timeline" quick filter to "event" resource type if no explicit timeline type exists
const quickFilters = [
	{ label: "All", value: "" },
	{ label: "Consent", value: "consent" },
	{ label: "Timeline", value: "event" }, 
	{ label: "Files", value: "file" },
	{ label: "Auth", value: "auth" },
] as const;

export function AuditLogPage(): JSX.Element {
	// Global State (Persisted)
	const { filters, setFilters, page, pageSize, setPage, setPageSize } =
		useAuditStore();

	// Local UI State
	const [activeQuickFilter, setActiveQuickFilter] = useState<string>(
		filters.resourceType || "",
	);
	const [showAdvancedFilters, setShowAdvancedFilters] = useState(false);
	const [selectedBatchId, setSelectedBatchId] = useState<string | null>(null);
	const [merkleTreeOpen, setMerkleTreeOpen] = useState(false);

	const queryParams = useMemo(() => {
		return {
			actor: filters.actor || undefined,
			resourceId: filters.resourceId || undefined,
			resourceType: (filters.resourceType || undefined) as AuditTargetType,
			action: (filters.action || undefined) as AuditAction,
			startTime: toIsoString(filters.startTime),
			endTime: toIsoString(filters.endTime),
			limit: pageSize,
			offset: (page - 1) * pageSize,
		};
	}, [filters, page, pageSize]);

	const auditQuery = useQuery({
		queryKey: ["audit-entries", queryParams],
		queryFn: () => queryAuditEntries(queryParams),
		placeholderData: (previousData) => previousData,
	});

	const batchesQuery = useQuery({
		queryKey: ["merkle-batches", { limit: 10, offset: 0 }],
		queryFn: async () => {
			const res = await listMerkleBatches({ limit: 10, offset: 0 });
			return res.batches;
		},
	});

	const selectedBatch = useMemo(() => {
		if (!selectedBatchId || !batchesQuery.data) return null;
		return batchesQuery.data.find((b) => b.id === selectedBatchId) || null;
	}, [selectedBatchId, batchesQuery.data]);

	const handleQuickFilter = (value: string) => {
		setActiveQuickFilter(value);
		setFilters({ resourceType: value as AuditTargetType });
		if (value) {
			setFilters({ action: "" }); // Clear action filter when switching types
		}
	};

	const handleFilterChange = (newFilters: AuditFilterState) => {
		setFilters(newFilters);
	};

	const entries = auditQuery.data?.entries ?? [];
	const totalCount = auditQuery.data?.total ?? 0;

	return (
		<div className="min-h-screen bg-background px-4 py-6 md:px-8 md:py-10">
			<div className="mx-auto flex w-full max-w-5xl flex-col gap-6">
				{/* Header */}
				<div className="flex flex-col gap-4">
					<Button
						variant="ghost"
						size="sm"
						asChild
						className="w-fit gap-2 text-muted-foreground hover:text-foreground"
						aria-label="Back to timeline"
					>
						<Link to="/">
							<ArrowLeft className="h-4 w-4" />
							Back to Timeline
						</Link>
					</Button>

					<div className="flex items-start justify-between gap-4">
						<div className="flex items-start gap-4">
							<div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 shrink-0">
								<FileText className="h-6 w-6 text-primary" />
							</div>
							<div className="flex flex-col gap-1">
								<h1 className="text-2xl font-bold tracking-tight text-foreground">
									Audit Log
								</h1>
								<p className="text-muted-foreground">
									Review the cryptographic audit trail for your account activity.
								</p>
							</div>
						</div>
					</div>
				</div>

				{/* Status Bar - Full Width, Prominent */}
				<AuditStatusBar />

				{/* Merkle Tree & Batch Management - Collapsible Section */}
				<div className="space-y-4">
					<Button
						variant="outline"
						onClick={() => setMerkleTreeOpen(!merkleTreeOpen)}
						className="w-full justify-between gap-2 border-dashed"
					>
						<span className="flex items-center gap-2">
							🌳 Merkle Tree & On-Chain Anchoring
						</span>
						<ChevronDown
							className={cn(
								"h-4 w-4 transition-transform duration-200",
								merkleTreeOpen && "rotate-180",
							)}
						/>
					</Button>

					{merkleTreeOpen && (
						<div className="space-y-4 animate-in fade-in slide-in-from-top-2 duration-200">
							{/* Batch management & verification */}
							{/* Manual manager removed as per requirement */}
                            <div className="flex justify-center">
							    <VerifyRootWidget />
                            </div>

							{/* Batch selector for tree view */}
							{batchesQuery.data && batchesQuery.data.length > 0 && (
								<div className="flex flex-wrap items-center gap-2">
									<span className="text-sm text-muted-foreground">
										View Merkle tree:
									</span>
									{batchesQuery.data.slice(0, 5).map((batch) => (
										<button
											key={batch.id}
											type="button"
											onClick={() =>
												setSelectedBatchId(
													selectedBatchId === batch.id ? null : batch.id,
												)
											}
											className={cn(
												"px-3 py-1.5 rounded-lg text-xs font-medium cursor-pointer",
												"border transition-all duration-200",
												"hover:border-primary/50",
												selectedBatchId === batch.id
													? "bg-primary/10 text-primary border-primary/50"
													: "bg-card border-border text-muted-foreground hover:text-foreground",
											)}
										>
											Batch #{batch.id}
											{batch.anchorStatus === "anchored" && (
												<span className="ml-1 text-success">✓</span>
											)}
										</button>
									))}
								</div>
							)}

							{/* Merkle tree viewer (if batch selected) */}
							{selectedBatch && (
								<div className="animate-in fade-in slide-in-from-top-2 duration-200">
									<MerkleTreeViewer batch={selectedBatch} />
								</div>
							)}
						</div>
					)}
				</div>

				{/* Quick filter chips & Advanced Filters */}
				<div className="flex flex-col gap-4">
                    <div className="flex flex-wrap items-center gap-2">
                        {quickFilters.map((filter) => (
                            <button
                                key={filter.value}
                                type="button"
                                onClick={() => handleQuickFilter(filter.value)}
                                className={cn(
                                    "px-3 py-1.5 rounded-full text-sm font-medium cursor-pointer",
                                    "border transition-all duration-200",
                                    "hover:border-primary/50",
                                    activeQuickFilter === filter.value
                                        ? "bg-primary text-primary-foreground border-primary"
                                        : "bg-card border-border text-muted-foreground hover:text-foreground",
                                )}
                            >
                                {filter.label}
                            </button>
                        ))}

                        <div className="flex-1" />

                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setShowAdvancedFilters(!showAdvancedFilters)}
                            className="gap-2"
                        >
                            <Filter className="h-4 w-4" />
                            {showAdvancedFilters ? "Hide Filters" : "Filters"}
                        </Button>

                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => auditQuery.refetch()}
                            disabled={auditQuery.isFetching}
                            className="gap-2"
                        >
                            <RefreshCw
                                className={cn("h-4 w-4", auditQuery.isFetching && "animate-spin")}
                            />
                            Refresh
                        </Button>
                    </div>

                    {/* Advanced filters (collapsible) */}
                    {showAdvancedFilters && (
                        <div className="animate-in fade-in slide-in-from-top-2 duration-200">
                            <AuditFilters
                                filters={filters}
                                onChange={handleFilterChange}
                                onReset={() => {
                                    setFilters(defaultFilters);
                                    setActiveQuickFilter("");
                                }}
                            />
                        </div>
                    )}
                </div>



				{/* Audit log entries */}
				<div className="space-y-4">
					<div className="rounded-md border bg-card">
						<AuditLogTable
							entries={entries}
							isLoading={auditQuery.isLoading}
							emptyState={
								<EmptyState
									icon={FileText}
									title="No audit entries found"
									description="Try adjusting your filters or upload some records to generate activity."
									action={
										<Button asChild variant="outline" size="sm" className="gap-2">
											<Link to="/">
												<Upload className="h-3 w-3" />
												Upload Record
											</Link>
										</Button>
									}
								/>
							}
						/>
					</div>

					<PaginationControls
						page={page}
						pageSize={pageSize}
						total={totalCount}
						onPageChange={setPage}
						onPageSizeChange={setPageSize}
						isLoading={auditQuery.isLoading}
					/>
				</div>
			</div>
		</div>
	);
}
