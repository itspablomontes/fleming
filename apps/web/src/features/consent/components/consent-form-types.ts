import type { ReactFormApi } from "@tanstack/react-form";
import type { ConsentRequestFormValues } from "./consent-request-wizard-types";

/**
 * Type for TanStack Form instances used in consent request components.
 *
 * TanStack Form's ReactFormApi requires 12 type parameters for validators and handlers.
 * Since these are complex and inferred by useForm, we use `any` for the validator types
 * as a pragmatic solution to allow form instances to be passed between components.
 */
export type ConsentForm = ReactFormApi<
	ConsentRequestFormValues,
	undefined,
	undefined,
	undefined,
	undefined,
	undefined,
	undefined,
	undefined,
	undefined,
	undefined,
	undefined,
	undefined
>;
