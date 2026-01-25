import type { JSX } from "react";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";

import type { ConsentForm } from "../consent-form-types";

interface DurationStepProps {
	form: ConsentForm;
}

const durationOptions = [
	{ value: "1", label: "1 day" },
	{ value: "7", label: "7 days" },
	{ value: "30", label: "30 days" },
	{ value: "90", label: "90 days" },
	{ value: "custom", label: "Custom" },
];

const presetValues = new Set([1, 7, 30, 90]);

export function DurationStep({ form }: DurationStepProps): JSX.Element {
	return (
		<form.Field name="durationDays">
			{(field) => {
				const durationValue = field.state.value;
				const mode =
					durationValue && presetValues.has(durationValue)
						? String(durationValue)
						: "custom";

				return (
					<div className="space-y-3">
						<Label htmlFor="duration-select">Access duration</Label>
						<Select
							value={mode}
							onValueChange={(value) => {
								if (value === "custom") {
									field.handleChange(undefined);
									return;
								}
								const parsed = Number.parseInt(value, 10);
								field.handleChange(Number.isNaN(parsed) ? undefined : parsed);
							}}
						>
							<SelectTrigger id="duration-select" className="w-full">
								<SelectValue placeholder="Select duration" />
							</SelectTrigger>
							<SelectContent>
								{durationOptions.map((option) => (
									<SelectItem key={option.value} value={option.value}>
										{option.label}
									</SelectItem>
								))}
							</SelectContent>
						</Select>

						{mode === "custom" && (
							<div className="space-y-2">
								<Label htmlFor="duration-custom">Custom days</Label>
								<Input
									id="duration-custom"
									type="number"
									min={1}
									max={365}
									value={field.state.value ?? ""}
									onChange={(event) => {
										const parsed = Number.parseInt(event.target.value, 10);
										field.handleChange(
											Number.isNaN(parsed) ? undefined : parsed,
										);
									}}
									placeholder="Enter number of days"
								/>
							</div>
						)}

						<p className="text-xs text-muted-foreground">
							Access automatically expires after the selected duration.
						</p>
					</div>
				);
			}}
		</form.Field>
	);
}
