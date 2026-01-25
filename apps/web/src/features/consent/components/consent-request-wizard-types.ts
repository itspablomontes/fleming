import type { ConsentPermission } from "../types";

export interface ConsentRequestFormValues {
	grantor: string;
	permissions: ConsentPermission[];
	durationDays?: number;
	reason: string;
}
