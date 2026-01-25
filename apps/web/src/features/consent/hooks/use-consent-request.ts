import { type UseMutationOptions, useMutation } from "@tanstack/react-query";

import { requestConsent } from "../api";
import type { ConsentGrant, ConsentRequestPayload } from "../types";

export const useConsentRequest = (
	options?: UseMutationOptions<ConsentGrant, Error, ConsentRequestPayload>,
) =>
	useMutation<ConsentGrant, Error, ConsentRequestPayload>({
		mutationFn: requestConsent,
		...options,
	});
