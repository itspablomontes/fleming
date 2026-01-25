import { useForm } from "@tanstack/react-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { JSX } from "react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { EthAddress } from "@/types/ethereum";

import { requestConsent } from "../api";
import type {
	ConsentGrant,
	ConsentPermission,
	ConsentRequestPayload,
} from "../types";
import type { ConsentForm } from "./consent-form-types";
import type { ConsentRequestFormValues } from "./consent-request-wizard-types";
import { StepIndicator } from "./step-indicator";
import { DurationStep } from "./wizard-steps/duration-step";
import { PatientStep } from "./wizard-steps/patient-step";
import { PermissionsStep } from "./wizard-steps/permissions-step";
import { ReasonStep } from "./wizard-steps/reason-step";
import { ReviewStep } from "./wizard-steps/review-step";

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

const stepLabels = [
	"Patient",
	"Permissions",
	"Duration",
	"Reason",
	"Review",
] as const;

interface ConsentRequestWizardProps {
	onSuccess?: (grant: ConsentGrant) => void;
}

export function ConsentRequestWizard({
	onSuccess,
}: ConsentRequestWizardProps): JSX.Element {
	const [currentStep, setCurrentStep] = useState(0);
	const queryClient = useQueryClient();

	const mutation = useMutation({
		mutationFn: (payload: ConsentRequestPayload) => requestConsent(payload),
		onSuccess: (grant) => {
			toast.success("Consent request sent.");
			onSuccess?.(grant);
			void queryClient.invalidateQueries({ queryKey: ["consent", "grants"] });
			form.reset();
			setCurrentStep(0);
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
			durationDays: 30 as number | undefined,
			reason: "",
		} satisfies Partial<ConsentRequestFormValues>,
		onSubmit: async ({ value }) => {
			const parsed = consentRequestSchema.safeParse(value);
			if (!parsed.success) {
				toast.error("Please fix the form errors before submitting.");
				return;
			}
			const payload: ConsentRequestPayload = {
				grantor: parsed.data.grantor as EthAddress,
				permissions: parsed.data.permissions as ConsentPermission[],
				reason: parsed.data.reason?.trim() || undefined,
				durationDays: parsed.data.durationDays,
			};
			await mutation.mutateAsync(payload);
		},
	});

	const values = form.state.values;
	const canProceed = useMemo(() => {
		switch (currentStep) {
			case 0:
				return consentRequestSchema.shape.grantor.safeParse(values.grantor)
					.success;
			case 1:
				return consentRequestSchema.shape.permissions.safeParse(
					values.permissions,
				).success;
			case 2:
				return consentRequestSchema.shape.durationDays.safeParse(
					values.durationDays,
				).success;
			case 3:
				return true;
			default:
				return form.state.canSubmit;
		}
	}, [currentStep, values, form.state.canSubmit]);

	const isLastStep = currentStep === stepLabels.length - 1;

	const handleNext = () => {
		if (!canProceed) {
			toast.error("Please complete this step before continuing.");
			return;
		}
		setCurrentStep((prev) => Math.min(prev + 1, stepLabels.length - 1));
	};

	const handleBack = () => {
		setCurrentStep((prev) => Math.max(prev - 1, 0));
	};

	return (
		<Card className="border-border bg-white dark:bg-gray-900">
			<CardHeader className="space-y-3">
				<CardTitle className="text-base font-semibold text-foreground">
					Request access
				</CardTitle>
				<p className="text-sm text-muted-foreground">
					Guide patients through a clear, secure consent request.
				</p>
				<StepIndicator steps={stepLabels} currentStep={currentStep} />
			</CardHeader>
			<CardContent>
				<form
					onSubmit={(event) => {
						event.preventDefault();
						event.stopPropagation();
						void form.handleSubmit();
					}}
					className="space-y-6"
				>
					<div aria-live="polite" className="sr-only">
						{`Step ${currentStep + 1} of ${stepLabels.length}: ${
							stepLabels[currentStep]
						}`}
					</div>

					{currentStep === 0 && <PatientStep form={form as ConsentForm} />}
					{currentStep === 1 && <PermissionsStep form={form as ConsentForm} />}
					{currentStep === 2 && <DurationStep form={form as ConsentForm} />}
					{currentStep === 3 && <ReasonStep form={form as ConsentForm} />}
					{currentStep === 4 && <ReviewStep values={values} />}

					<div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
						<Button
							type="button"
							variant="outline"
							disabled={currentStep === 0 || mutation.isPending}
							onClick={handleBack}
						>
							Back
						</Button>

						{currentStep === 3 && (
							<Button
								type="button"
								variant="ghost"
								onClick={handleNext}
								disabled={mutation.isPending}
							>
								Skip
							</Button>
						)}

						{isLastStep ? (
							<form.Subscribe
								selector={(state) => [state.canSubmit, state.isSubmitting]}
							>
								{([canSubmit, isSubmitting]) => (
									<Button type="submit" disabled={!canSubmit || isSubmitting}>
										{isSubmitting ? "Sending request..." : "Submit request"}
									</Button>
								)}
							</form.Subscribe>
						) : (
							<Button
								type="button"
								onClick={handleNext}
								disabled={mutation.isPending}
							>
								Next
							</Button>
						)}
					</div>
				</form>
			</CardContent>
		</Card>
	);
}
