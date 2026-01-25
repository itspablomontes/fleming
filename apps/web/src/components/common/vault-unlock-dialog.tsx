import { Lock } from "lucide-react";
import React from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useVault } from "@/features/auth/contexts/vault-context";

interface VaultUnlockDialogProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

export function VaultUnlockDialog({
  isOpen,
  onOpenChange,
  onSuccess,
}: VaultUnlockDialogProps) {
  const { unlockVault, isUnlocking } = useVault();
  const [error, setError] = React.useState<Error | null>(null);

  const handleUnlock = async () => {
    setError(null);
    try {
      await unlockVault();
      onSuccess?.();
      onOpenChange(false);
    } catch (error) {
       setError(error instanceof Error ? error : new Error("Failed to unlock vault"));
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md bg-white dark:bg-gray-950 border-cyan-200 dark:border-cyan-900">
        <DialogHeader>
          <div className="mx-auto bg-cyan-100 dark:bg-cyan-900/30 p-3 rounded-full mb-4">
            <Lock className="h-6 w-6 text-cyan-600 dark:text-cyan-400" />
          </div>
          <DialogTitle className="text-center text-xl text-cyan-900 dark:text-cyan-50">
            Unlock Your Medical Vault
          </DialogTitle>
          <DialogDescription className="text-center text-gray-500 dark:text-gray-400">
            Your medical data is end-to-end encrypted. We need your signature to derived the keys to unlock it.
            <br />
            <br />
            <span className="text-xs">
              We never see your private keys. The signature helps us mathematically regenerate your secure session key.
            </span>
          </DialogDescription>
        </DialogHeader>
        <div className="px-6 pb-2">
            {error && (
            <div className="mb-4 p-3 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded-md text-sm border border-red-200 dark:border-red-900">
                {error.message}
            </div>
            )}
        </div>
        <DialogFooter className="sm:justify-center">
          <Button
            onClick={handleUnlock}
            disabled={isUnlocking}
            className="w-full sm:w-auto bg-cyan-600 hover:bg-cyan-700 text-white"
          >
            {isUnlocking ? "Verifying Signature..." : "Sign to Unlock"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
