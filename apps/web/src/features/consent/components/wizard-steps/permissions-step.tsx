import type { JSX } from "react";

import { Label } from "@/components/ui/label";

import type { ConsentPermission } from "../../types";
import { ConsentPermission as PermissionValue } from "../../types";
import type { ConsentForm } from "../consent-form-types";

interface PermissionsStepProps {
	form: ConsentForm;
}

const permissionOptions: Array<{
	value: ConsentPermission;
	label: string;
	description: string;
}> = [
	{
		value: PermissionValue.Read,
		label: "Read",
		description: "View timeline events and encrypted records.",
	},
	{
		value: PermissionValue.Write,
		label: "Write",
		description: "Add new records to the patient timeline.",
	},
	{
		value: PermissionValue.Share,
		label: "Share",
		description: "Re-share records with other providers.",
	},
];

export function PermissionsStep({ form }: PermissionsStepProps): JSX.Element {
	return (
		<form.Field
			name="permissions"
			validators={{
				onChange: ({ value }: { value: ConsentPermission[] }) =>
					value.length > 0 ? undefined : "Select at least one permission",
			}}
		>
			{(field) => (
				<div className="space-y-3">
					<Label>Requested permissions</Label>
					<div className="space-y-3">
						{permissionOptions.map((option) => {
							const isChecked = field.state.value.includes(option.value);
							return (
								<label
									key={option.value}
									className="flex cursor-pointer items-start gap-3 rounded-lg border border-border p-3 text-sm transition-colors hover:bg-muted"
								>
									<input
										type="checkbox"
										className="mt-1 h-4 w-4 rounded border-gray-300 accent-cyan-600 focus-visible:ring-2 focus-visible:ring-cyan-500 focus-visible:ring-offset-2"
										checked={isChecked}
										onChange={() => {
											const next = isChecked
												? field.state.value.filter(
														(value: ConsentPermission) =>
															value !== option.value,
													)
												: [...field.state.value, option.value];
											field.handleChange(next);
										}}
									/>
									<div className="space-y-1">
										<p className="font-medium text-foreground">
											{option.label}
										</p>
										<p className="text-xs text-muted-foreground">
											{option.description}
										</p>
									</div>
								</label>
							);
						})}
					</div>
					{field.state.meta.errors.length > 0 && (
						<p
							role="alert"
							aria-live="polite"
							className="text-xs text-destructive"
						>
							{field.state.meta.errors[0]}
						</p>
					)}
				</div>
			)}
		</form.Field>
	);
}
