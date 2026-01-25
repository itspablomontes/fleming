import type { JSX } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import type { MerkleProof } from "../types";

interface MerkleProofViewerProps {
	root?: string;
	proof?: MerkleProof;
	isValid?: boolean;
}

export function MerkleProofViewer({
	root,
	proof,
	isValid,
}: MerkleProofViewerProps): JSX.Element {
	if (!proof) {
		return (
			<Card className="bg-white dark:bg-gray-900">
				<CardHeader>
					<CardTitle className="text-sm font-semibold text-gray-900 dark:text-gray-100">
						Merkle Proof
					</CardTitle>
				</CardHeader>
				<CardContent className="text-xs text-muted-foreground">
					No proof loaded yet.
				</CardContent>
			</Card>
		);
	}

	return (
		<Card className="bg-white dark:bg-gray-900">
			<CardHeader>
				<CardTitle className="text-sm font-semibold text-gray-900 dark:text-gray-100">
					Merkle Proof {isValid !== undefined && (isValid ? "✓" : "✗")}
				</CardTitle>
			</CardHeader>
			<CardContent className="space-y-2 text-xs text-muted-foreground">
				{root && (
					<p className="truncate">
						<span className="font-medium text-gray-700 dark:text-gray-200">
							Root:
						</span>{" "}
						{root}
					</p>
				)}
				<p className="truncate">
					<span className="font-medium text-gray-700 dark:text-gray-200">
						Entry:
					</span>{" "}
					{proof.entryHash}
				</p>
				<Separator />
				<div className="space-y-2">
					{proof.steps.map((step, index) => (
						<p key={`${step.hash}-${index}`} className="truncate">
							<span className="font-medium text-gray-700 dark:text-gray-200">
								{step.isLeft ? "Left" : "Right"}:
							</span>{" "}
							{step.hash}
						</p>
					))}
				</div>
			</CardContent>
		</Card>
	);
}
