import {
	ChevronFirst,
	ChevronLast,
	ChevronLeft,
	ChevronRight,
} from "lucide-react";
import type { JSX } from "react";

import { Button } from "@/components/ui/button";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";

interface PaginationControlsProps {
	page: number;
	pageSize: number;
	total: number;
	onPageChange: (page: number) => void;
	onPageSizeChange: (pageSize: number) => void;
	isLoading?: boolean;
}

export function PaginationControls({
	page,
	pageSize,
	total,
	onPageChange,
	onPageSizeChange,
	isLoading,
}: PaginationControlsProps): JSX.Element {
	const totalPages = Math.ceil(total / pageSize);
	const startItem = (page - 1) * pageSize + 1;
	// Fix endItem logic to correct display on last page
	const endItem = Math.min(page * pageSize, total);

	if (total === 0) return <div />;

	return (
		<div className="flex flex-col items-center justify-between gap-4 px-2 py-4 sm:flex-row">
			<div className="flex items-center gap-2 text-sm text-muted-foreground">
				<span>Rows per page</span>
				<Select
					value={String(pageSize)}
					onValueChange={(v) => onPageSizeChange(Number(v))}
					disabled={isLoading}
				>
					<SelectTrigger className="h-8 w-[70px]">
						<SelectValue placeholder={pageSize} />
					</SelectTrigger>
					<SelectContent side="top">
						{[10, 20, 50, 100].map((size) => (
							<SelectItem key={size} value={String(size)}>
								{size}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
				<span className="hidden sm:inline-block ml-4">
					Showing {startItem}-{endItem} of {total}
				</span>
			</div>

			<div className="flex items-center space-x-2">
				<Button
					variant="outline"
					size="icon"
					className="h-8 w-8"
					onClick={() => onPageChange(1)}
					disabled={page === 1 || isLoading}
					aria-label="First page"
				>
					<ChevronFirst className="h-4 w-4" />
				</Button>
				<Button
					variant="outline"
					size="icon"
					className="h-8 w-8"
					onClick={() => onPageChange(page - 1)}
					disabled={page === 1 || isLoading}
					aria-label="Previous page"
				>
					<ChevronLeft className="h-4 w-4" />
				</Button>

				<div className="flex w-[100px] items-center justify-center text-sm font-medium">
					Page {page} of {totalPages}
				</div>

				<Button
					variant="outline"
					size="icon"
					className="h-8 w-8"
					onClick={() => onPageChange(page + 1)}
					disabled={page >= totalPages || isLoading}
					aria-label="Next page"
				>
					<ChevronRight className="h-4 w-4" />
				</Button>
				<Button
					variant="outline"
					size="icon"
					className="h-8 w-8"
					onClick={() => onPageChange(totalPages)}
					disabled={page >= totalPages || isLoading}
					aria-label="Last page"
				>
					<ChevronLast className="h-4 w-4" />
				</Button>
			</div>
		</div>
	);
}
