import { useMutation, useQuery } from "@tanstack/react-query";
import type { JSX } from "react";
import { useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";

import {
	anchorMerkleBatch,
	buildMerkleTree,
	listMerkleBatches,
	type MerkleBatch,
} from "../api";

function formatDateTime(value?: string): string {
	if (!value) return "";
	const d = new Date(value);
	if (Number.isNaN(d.getTime())) return value;
	return d.toLocaleString();
}

function shortHash(value?: string): string {
	if (!value) return "";
	if (value.length <= 16) return value;
	return `${value.slice(0, 10)}…${value.slice(-6)}`;
}

export function MerkleBatchManager(): JSX.Element {
	const [startTime, setStartTime] = useState<string>("");
	const [endTime, setEndTime] = useState<string>("");

	const batchesQuery = useQuery({
		queryKey: ["merkle-batches", { limit: 25, offset: 0 }],
		queryFn: async () => {
			const res = await listMerkleBatches({ limit: 25, offset: 0 });
			return res.batches;
		},
	});

	const buildMutation = useMutation({
		mutationFn: async () => {
			return buildMerkleTree({
				startTime: startTime ? new Date(startTime) : undefined,
				endTime: endTime ? new Date(endTime) : undefined,
			});
		},
		onSuccess: async () => {
			await batchesQuery.refetch();
		},
	});

	const anchorMutation = useMutation({
		mutationFn: async (batchId: string) => {
			return anchorMerkleBatch(batchId);
		},
		onSuccess: async () => {
			await batchesQuery.refetch();
		},
	});

	const buildResult = buildMutation.data?.batch;

	const canBuild = useMemo(() => {
		// Allow empty (meaning “all time”), but avoid invalid dates from datetime-local input.
		if (startTime && Number.isNaN(new Date(startTime).getTime())) return false;
		if (endTime && Number.isNaN(new Date(endTime).getTime())) return false;
		return true;
	}, [startTime, endTime]);

	const renderRow = (batch: MerkleBatch) => {
		const anchored = batch.anchorStatus === "anchored";
		const isAnchoring =
			anchorMutation.isPending && anchorMutation.variables === batch.id;

		return (
			<div key={batch.id} className="flex flex-col gap-2 rounded-md border p-3">
				<div className="flex items-start justify-between gap-4">
					<div className="min-w-0">
						<p className="text-sm font-medium text-gray-900 dark:text-gray-100">
							Batch {batch.id}
						</p>
						<p className="text-xs text-muted-foreground">
							Created {formatDateTime(batch.createdAt)} • Entries{" "}
							{batch.entryCount}
						</p>
					</div>
					<Button
						variant="outline"
						size="sm"
						disabled={anchored || anchorMutation.isPending}
						onClick={() => anchorMutation.mutate(batch.id)}
					>
						{anchored ? "Anchored" : isAnchoring ? "Anchoring..." : "Anchor"}
					</Button>
				</div>

				<p className="text-xs text-muted-foreground truncate">
					Root: {shortHash(batch.rootHash)}
				</p>
				<p className="text-xs text-muted-foreground">
					Status: {batch.anchorStatus ?? "pending"}
					{batch.anchoredAt
						? ` • Anchored ${formatDateTime(batch.anchoredAt)}`
						: ""}
				</p>
				{batch.anchorTxHash ? (
					<p className="text-xs text-muted-foreground truncate">
						Tx: {batch.anchorTxHash}
					</p>
				) : null}
				{batch.anchorError ? (
					<p className="text-xs text-destructive">Error: {batch.anchorError}</p>
				) : null}
			</div>
		);
	};

	return (
		<Card className="bg-white dark:bg-gray-900">
			<CardHeader className="flex flex-col gap-1">
				<CardTitle className="text-sm font-semibold text-gray-900 dark:text-gray-100">
					Merkle batches (on-chain anchoring)
				</CardTitle>
				<p className="text-xs text-muted-foreground">
					Build an actor-scoped Merkle root for a time window, then anchor it on
					the local chain.
				</p>
			</CardHeader>
			<CardContent className="space-y-4">
				<div className="grid gap-3 md:grid-cols-3">
					<div className="space-y-2">
						<p className="text-xs font-medium text-gray-700 dark:text-gray-200">
							Start time (optional)
						</p>
						<Input
							type="datetime-local"
							value={startTime}
							onChange={(e) => setStartTime(e.target.value)}
						/>
					</div>
					<div className="space-y-2">
						<p className="text-xs font-medium text-gray-700 dark:text-gray-200">
							End time (optional)
						</p>
						<Input
							type="datetime-local"
							value={endTime}
							onChange={(e) => setEndTime(e.target.value)}
						/>
					</div>
					<div className="flex items-end">
						<Button
							className="w-full"
							disabled={!canBuild || buildMutation.isPending}
							onClick={() => buildMutation.mutate()}
						>
							{buildMutation.isPending ? "Building..." : "Build batch"}
						</Button>
					</div>
				</div>

				{buildMutation.isError ? (
					<p className="text-xs text-destructive">
						Failed to build Merkle batch.
					</p>
				) : null}

				{buildResult ? (
					<div className="rounded-md border p-3">
						<p className="text-xs text-muted-foreground">
							Latest build result: batch {buildResult.id} • root{" "}
							{shortHash(buildResult.rootHash)}
						</p>
					</div>
				) : null}

				<Separator />

				<div className="flex items-center justify-between">
					<p className="text-sm font-medium text-gray-900 dark:text-gray-100">
						Recent batches
					</p>
					<Button
						variant="outline"
						size="sm"
						onClick={() => batchesQuery.refetch()}
						disabled={batchesQuery.isFetching}
					>
						{batchesQuery.isFetching ? "Refreshing..." : "Refresh"}
					</Button>
				</div>

				{batchesQuery.isLoading ? (
					<p className="text-xs text-muted-foreground">Loading batches…</p>
				) : batchesQuery.isError ? (
					<p className="text-xs text-destructive">Failed to load batches.</p>
				) : (batchesQuery.data?.length ?? 0) === 0 ? (
					<p className="text-xs text-muted-foreground">No batches yet.</p>
				) : (
					<div className="space-y-2">{batchesQuery.data?.map(renderRow)}</div>
				)}
			</CardContent>
		</Card>
	);
}
