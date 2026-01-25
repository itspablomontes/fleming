import type { JSX } from "react";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

import type { ConsentForm } from "../consent-form-types";

interface PatientStepProps {
	form: ConsentForm;
}

const addressRegex = /^0x[a-fA-F0-9]{40}$/;

export function PatientStep({ form }: PatientStepProps): JSX.Element {
	return (
		<form.Field
			name="grantor"
			validators={{
				onChange: ({ value }: { value: string }) =>
					addressRegex.test(value.trim())
						? undefined
						: "Enter a valid Ethereum address",
			}}
		>
			{(field) => (
				<div className="space-y-3">
					<div className="space-y-2">
						<Label htmlFor={field.name}>Patient wallet address</Label>
						<Input
							id={field.name}
							value={field.state.value}
							onBlur={field.handleBlur}
							onChange={(event) => field.handleChange(event.target.value)}
							placeholder="0x..."
							aria-invalid={Boolean(field.state.meta.errors.length)}
							aria-describedby={`${field.name}-error`}
						/>
						{field.state.meta.errors.length > 0 && (
							<p
								id={`${field.name}-error`}
								role="alert"
								aria-live="polite"
								className="text-xs text-destructive"
							>
								{field.state.meta.errors[0]}
							</p>
						)}
					</div>
					<div
						className={cn(
							"rounded-md border border-border bg-muted/40 p-3 text-xs text-muted-foreground",
						)}
					>
						Use the patient’s wallet address to request access. This address is
						never stored in plaintext on the server.
					</div>
				</div>
			)}
		</form.Field>
	);
}
