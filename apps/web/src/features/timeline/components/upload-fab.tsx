import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";

interface UploadFABProps {
	onClick: () => void;
}

export function UploadFAB({ onClick }: UploadFABProps) {
	return (
		<Button
			onClick={onClick}
			size="lg"
			className="fixed bottom-6 left-6 h-14 w-14 rounded-full shadow-lg bg-emerald-600 hover:bg-emerald-700 text-white transition-all hover:scale-110 focus:ring-4 focus:ring-emerald-500/50 z-50"
			aria-label="Upload document"
		>
			<Plus className="h-6 w-6" />
		</Button>
	);
}
