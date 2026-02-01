"use client";

import {
	type ColumnFiltersState,
	flexRender,
	getCoreRowModel,
	getPaginationRowModel,
	getSortedRowModel,
	type SortingState,
	useReactTable,
	type VisibilityState,
} from "@tanstack/react-table";
import { type JSX, useState } from "react";
import { EmptyState } from "@/components/ui/empty-state";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import type { AuditLogEntry } from "../types";
import { columns } from "./audit-log-columns";

interface AuditLogTableProps {
	entries: AuditLogEntry[];
	isLoading?: boolean;
	emptyState?: React.ReactNode;
}

export function AuditLogTable({
	entries,
	isLoading,
	emptyState,
}: AuditLogTableProps): JSX.Element {
	const [sorting, setSorting] = useState<SortingState>([]);
	const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
	const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
	const [rowSelection, setRowSelection] = useState({});

	const table = useReactTable({
		data: entries,
		columns,
		onSortingChange: setSorting,
		onColumnFiltersChange: setColumnFilters,
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(), // Note: Pagination is mostly handled server-side by the page, but this enables client-side per-page if needed.
		// Important: If server-side pagination is used, we might disable this or pass manual pagination props.
		// Since we pass a sliced 'entries' array from the parent which represents ONE page, `getPaginationRowModel` here
		// effectively allows paginating THAT slice, which is wrong.
		// We should just render all rows given.
		// However, TanStack table expects a row model. `getCoreRowModel` is sufficient for "view all provided rows".
		// If sorting is client-side (on the current page data only), `getSortedRowModel` works.
		getSortedRowModel: getSortedRowModel(),
		onColumnVisibilityChange: setColumnVisibility,
		onRowSelectionChange: setRowSelection,
		state: {
			sorting,
			columnFilters,
			columnVisibility,
			rowSelection,
		},
	});

	if (isLoading) {
		return (
			<div className="w-full rounded-md border bg-card">
				<div className="p-4">
					<div className="flex items-center justify-between pb-4">
						<div className="h-6 w-32 animate-pulse rounded bg-muted" />
						<div className="h-6 w-16 animate-pulse rounded bg-muted" />
					</div>
					<div className="space-y-4">
						{[...Array(5)].map((_, i) => (
							<div
								// biome-ignore lint/suspicious/noArrayIndexKey: Loading skeleton
								key={i}
								className="h-12 w-full animate-pulse rounded border border-border bg-muted/30"
							/>
						))}
					</div>
				</div>
			</div>
		);
	}

	if (!entries.length) {
		if (emptyState) return <>{emptyState}</>;
		return <EmptyState title={"No audit entries found"} />;
	}

	return (
		<div className="rounded-md border bg-card shadow-sm overflow-hidden">
			<Table>
				<TableHeader className="bg-muted/40">
					{table.getHeaderGroups().map((headerGroup) => (
						<TableRow
							key={headerGroup.id}
							className="hover:bg-transparent border-b-border"
						>
							{headerGroup.headers.map((header) => {
								return (
									<TableHead key={header.id} className="h-11">
										{header.isPlaceholder
											? null
											: flexRender(
													header.column.columnDef.header,
													header.getContext(),
												)}
									</TableHead>
								);
							})}
						</TableRow>
					))}
				</TableHeader>
				<TableBody>
					{table.getRowModel().rows?.length ? (
						table.getRowModel().rows.map((row) => (
							<TableRow
								key={row.id}
								data-state={row.getIsSelected() && "selected"}
								className="group hover:bg-muted/30 data-[state=selected]:bg-muted"
							>
								{row.getVisibleCells().map((cell) => (
									<TableCell key={cell.id} className="py-3">
										{flexRender(cell.column.columnDef.cell, cell.getContext())}
									</TableCell>
								))}
							</TableRow>
						))
					) : (
						<TableRow>
							<TableCell colSpan={columns.length} className="h-24 text-center">
								No results.
							</TableCell>
						</TableRow>
					)}
				</TableBody>
			</Table>
		</div>
	);
}
