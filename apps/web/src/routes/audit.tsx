import { createFileRoute } from "@tanstack/react-router";
import { AuditLogPage } from "@/features/audit/pages/audit-log-page";

export const Route = createFileRoute("/audit")({
	component: AuditLogPage,
});
