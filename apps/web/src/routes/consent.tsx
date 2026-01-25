import { createFileRoute } from "@tanstack/react-router";
import { ConsentPage } from "@/features/consent/pages/consent-page";

export const Route = createFileRoute("/consent")({
	component: ConsentPage,
});
