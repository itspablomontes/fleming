import { useForm } from "@tanstack/react-form";
import { useMutation } from "@tanstack/react-query";
import type { JSX } from "react";
import { useState } from "react";
import { toast } from "sonner";
import { z } from "zod";
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
import { cn } from "@/lib/utils";
import type { EthAddress } from "@/types/ethereum";

import { requestConsent } from "../api";
import {
	type ConsentGrant,
	ConsentPermission,
	type ConsentRequestPayload,
} from "../types";

const permissionOptions: Array<{
	value: ConsentPermission;
	label: string;
	description: string;
}> = [
	{
		value: ConsentPermission.Read,
		label: "Read",
		description: "View the timeline and encrypted records.",
	},
	{
		value: ConsentPermission.Write,
		label: "Write",
		description: "Add new events to the timeline.",
	},
	{
		value: ConsentPermission.Share,
		label: "Share",
		description: "Re-share records with other providers.",
	},
];

const durationOptions = [
	{ value: "1", label: "1 day" },
	{ value: "7", label: "7 days" },
	{ value: "30", label: "30 days" },
	{ value: "90", label: "90 days" },
	{ value: "custom", label: "Custom" },
];

const consentRequestSchema = z.object({
	grantor: z
		.string()
		.trim()
		.regex(/^0x[a-fA-F0-9]{40}$/, "Enter a valid Ethereum address"),
	permissions: z
		.array(z.enum(["read", "write", "share"]))
		.min(1, "Select at least one permission"),
	reason: z.string().trim().max(500, "Reason is too long").optional(),
	durationDays: z
		.number()
		.int("Use a whole number of days")
		.positive("Duration must be greater than 0")
		.max(365, "Duration cannot exceed 365 days")
		.optional(),
});

type ConsentRequestFormValues = z.infer<typeof consentRequestSchema>;

interface ConsentRequestFormProps {
	onSuccess?: (grant: ConsentGrant) => void;
	className?: string;
}

export function ConsentRequestForm({
	onSuccess,
	className,
}: ConsentRequestFormProps): JSX.Element {
	const [durationMode, setDurationMode] = useState<string>("30");

	const mutation = useMutation({
		mutationFn: (payload: ConsentRequestPayload) => requestConsent(payload),
		onSuccess: (grant) => {
			toast.success("Consent request sent.");
			onSuccess?.(grant);
		},
		onError: (error) => {
			const message =
				error instanceof Error ? error.message : "Failed to request consent";
			toast.error(message);
		},
	});

	const form = useForm({
		defaultValues: {
			grantor: "",
			permissions: [] as ConsentPermission[],
			reason: "",
			durationDays: 30 as number | undefined,
		} satisfies Partial<ConsentRequestFormValues>,
		onSubmit: async ({ value }) => {
			const parsed = consentRequestSchema.safeParse(value);
			if (!parsed.success) {
				toast.error("Please fix the form errors before submitting.");
				return;
			}
			const payload: ConsentRequestPayload = {
				grantor: parsed.data.grantor as EthAddress,
				permissions: parsed.data.permissions,
				reason: parsed.data.reason?.trim() || undefined,
				durationDays: parsed.data.durationDays,
			};
			await mutation.mutateAsync(payload);
		},
	});

	return (
		<form
			onSubmit={(event) => {
				event.preventDefault();
				event.stopPropagation();
				void form.handleSubmit();
			}}
			className={cn("space-y-6", className)}
		>
			<form.Field
				name="grantor"
				validators={{
					onChange: ({ value }) =>
						consentRequestSchema.shape.grantor.safeParse(value).success
							? undefined
							: "Enter a valid Ethereum address",
				}}
			>
				{(field) => (
					<div className="space-y-2">
						<Label htmlFor={field.name}>Patient address</Label>
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
				)}
			</form.Field>

			<form.Field
				name="permissions"
				validators={{
					onChange: ({ value }: { value: ConsentPermission[] }) =>
						consentRequestSchema.shape.permissions.safeParse(value).success
							? undefined
							: "Select at least one permission",
				}}
			>
				{(field) => {
					const permissions = field.state.value as ConsentPermission[];
					return (
						<div className="space-y-2">
							<Label>Permissions</Label>
							<div className="space-y-3">
								{permissionOptions.map((option) => {
									const isChecked = permissions.includes(option.value);
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
														? permissions.filter(
																(value: ConsentPermission) =>
																	value !== option.value,
															)
														: [...permissions, option.value];
													field.handleChange(next as ConsentPermission[]);
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
					);
				}}
			</form.Field>

			<form.Field name="durationDays">
				{(field) => (
					<div className="space-y-2">
						<Label htmlFor="duration-select">Access duration</Label>
						<Select
							value={durationMode}
							onValueChange={(value) => {
								setDurationMode(value);
								if (value === "custom") {
									field.handleChange(undefined as number | undefined);
									return;
								}
								const parsed = Number.parseInt(value, 10);
								field.handleChange(
									(Number.isNaN(parsed) ? undefined : parsed) as
										| number
										| undefined,
								);
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
						{durationMode === "custom" && (
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
											(Number.isNaN(parsed) ? undefined : parsed) as
												| number
												| undefined,
										);
									}}
									placeholder="Enter number of days"
								/>
							</div>
						)}
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

			<form.Field name="reason">
				{(field) => (
					<div className="space-y-2">
						<Label htmlFor="reason">Reason (optional)</Label>
						<textarea
							id="reason"
							value={field.state.value}
							onChange={(event) => field.handleChange(event.target.value)}
							onBlur={field.handleBlur}
							placeholder="Explain why you need access"
							className="min-h-[96px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-xs focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-500 focus-visible:ring-offset-2 dark:focus-visible:ring-cyan-400"
						/>
					</div>
				)}
			</form.Field>

			<form.Subscribe
				selector={(state) => [state.canSubmit, state.isSubmitting]}
			>
				{([canSubmit, isSubmitting]) => (
					<div className="flex items-center justify-between gap-3">
						<Button
							type="submit"
							disabled={!canSubmit || isSubmitting}
							className="ml-auto"
						>
							{isSubmitting ? "Sending request..." : "Request access"}
						</Button>
					</div>
				)}
			</form.Subscribe>
		</form>
	);
}
