import { FileUp } from "lucide-react";
import { Button } from "@/components/ui/button";

interface TimelineEmptyStateProps {
	onUploadClick: () => void;
}

export function TimelineEmptyState({ onUploadClick }: TimelineEmptyStateProps) {
	return (
		<div className="flex flex-col items-center justify-center min-h-[60vh] px-4 text-center">
			{/* Illustration placeholder - can be replaced with actual SVG */}
			<div className="mb-8 w-64 h-64 rounded-full bg-cyan-50 dark:bg-cyan-950/20 flex items-center justify-center">
				<FileUp className="w-32 h-32 text-cyan-600 dark:text-cyan-400" />
			</div>

			<h2 className="text-2xl font-semibold text-cyan-900 dark:text-cyan-50 mb-3">
				No medical records yet
			</h2>

			<p className="text-cyan-700 dark:text-cyan-300 max-w-md mb-8">
				Upload your first document to start building your health timeline. All
				files are encrypted and only you control access.
			</p>

			<Button
				onClick={onUploadClick}
				size="lg"
				className="bg-emerald-600 hover:bg-emerald-700 text-white"
			>
				<FileUp className="mr-2 h-5 w-5" />
				Upload Document
			</Button>
		</div>
	);
}
