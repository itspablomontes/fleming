import type { JSX } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import {
	AuditAction,
	type AuditAction as AuditActionType,
	AuditTargetType,
	type AuditTargetType as AuditTargetTypeType,
} from "../types";

export interface AuditFilterState {
	actor: string;
	resourceId: string;
	resourceType: AuditTargetTypeType | "";
	action: AuditActionType | "";
	startTime: string;
	endTime: string;
}

interface AuditFiltersProps {
	filters: AuditFilterState;
	onChange: (filters: AuditFilterState) => void;
	onReset?: () => void;
}

const actionOptions = Object.values(AuditAction);
const resourceOptions = Object.values(AuditTargetType);

export function AuditFilters({
	filters,
	onChange,
	onReset,
}: AuditFiltersProps): JSX.Element {
	const update = <K extends keyof AuditFilterState>(
		key: K,
		value: AuditFilterState[K],
	) => {
		onChange({ ...filters, [key]: value });
	};

	return (
		<div className="grid gap-4 rounded-lg border border-gray-200 bg-white p-4 dark:border-gray-800 dark:bg-gray-900 md:grid-cols-2 lg:grid-cols-3">
			<div className="space-y-2">
				<Label htmlFor="audit-actor">Actor</Label>
				<Input
					id="audit-actor"
					value={filters.actor}
					onChange={(event) => update("actor", event.target.value)}
					placeholder="0x..."
				/>
			</div>
			<div className="space-y-2">
				<Label htmlFor="audit-resource">Resource ID</Label>
				<Input
					id="audit-resource"
					value={filters.resourceId}
					onChange={(event) => update("resourceId", event.target.value)}
					placeholder="resource id"
				/>
			</div>
			<div className="space-y-2">
				<Label>Resource Type</Label>
				<Select
					value={filters.resourceType || "all"}
					onValueChange={(value) =>
						update(
							"resourceType",
							value === "all" ? "" : (value as AuditTargetTypeType),
						)
					}
				>
					<SelectTrigger className="w-full">
						<SelectValue placeholder="All types" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="all">All</SelectItem>
						{resourceOptions.map((option) => (
							<SelectItem key={option} value={option}>
								{option}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
			</div>
			<div className="space-y-2">
				<Label>Action</Label>
				<Select
					value={filters.action || "all"}
					onValueChange={(value) =>
						update("action", value === "all" ? "" : (value as AuditActionType))
					}
				>
					<SelectTrigger className="w-full">
						<SelectValue placeholder="All actions" />
					</SelectTrigger>
					<SelectContent>
						<SelectItem value="all">All</SelectItem>
						{actionOptions.map((option) => (
							<SelectItem key={option} value={option}>
								{option}
							</SelectItem>
						))}
					</SelectContent>
				</Select>
			</div>
			<div className="space-y-2">
				<Label htmlFor="audit-start">Start Time</Label>
				<Input
					id="audit-start"
					type="datetime-local"
					value={filters.startTime}
					onChange={(event) => update("startTime", event.target.value)}
				/>
			</div>
			<div className="space-y-2">
				<Label htmlFor="audit-end">End Time</Label>
				<Input
					id="audit-end"
					type="datetime-local"
					value={filters.endTime}
					onChange={(event) => update("endTime", event.target.value)}
				/>
			</div>
			{onReset && (
				<div className="flex items-end">
					<Button variant="outline" onClick={onReset}>
						Reset filters
					</Button>
				</div>
			)}
		</div>
	);
}
