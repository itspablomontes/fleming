import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { useConnection, useSignMessage } from "wagmi";
import { deriveMasterKey } from "@/lib/crypto/keys";

// User salt is fetched from backend profile

interface VaultContextType {
  isUnlocked: boolean;
  isUnlocking: boolean;
  masterKey: CryptoKey | null;
  unlockVault: () => Promise<void>;
  lockVault: () => void;
}

const VaultContext = createContext<VaultContextType | undefined>(undefined);

export function VaultProvider({ children }: { children: React.ReactNode }) {
  const { address, status } = useConnection();
  const { mutateAsync: signMessageAsync } = useSignMessage();

  const [isUnlocked, setIsUnlocked] = useState(false);
  const [isUnlocking, setIsUnlocking] = useState(false);
  const [masterKey, setMasterKey] = useState<CryptoKey | null>(null);

  const lockVault = useCallback(() => {
    setIsUnlocked(false);
    setMasterKey(null);
    setIsUnlocking(false);
  }, []);

  // Auto-lock when wallet disconnects or changes
  useEffect(() => {
		if (status !== "connected" || !address) {
      lockVault();
    }
	}, [status, address, lockVault]);

  const unlockVault = useCallback(async () => {
		if (!address) return;
    setIsUnlocking(true);

    try {
      // 1. Request signature from wallet
      const signature = await signMessageAsync({
				message: `Unlock Fleming Vault for ${address}\n\nSign this message to derive your encryption keys.`,
      });

      // 2. Fetch User Salt from Backend
      // The salt is permanently bound to the user's account in the backend DB.
      const response = await fetch("/api/auth/me");
      if (!response.ok) throw new Error("Failed to fetch user profile");
      
      const data = await response.json();
      if (!data.encryptionSalt) {
        throw new Error("User has no encryption salt. Please contact support.");
      }

      // 3. Derive Master Key
      const saltBuffer = new TextEncoder().encode(data.encryptionSalt);
      const key = await deriveMasterKey(signature, saltBuffer);

      setMasterKey(key);
      setIsUnlocked(true);
    } catch (error) {
      console.error("Failed to unlock vault:", error);
      throw error;
    } finally {
      setIsUnlocking(false);
    }
	}, [address, signMessageAsync]);

  return (
    <VaultContext.Provider
      value={{
        isUnlocked,
        isUnlocking,
        masterKey,
        unlockVault,
        lockVault,
      }}
    >
      {children}
    </VaultContext.Provider>
  );
}

export function useVault() {
  const context = useContext(VaultContext);
  if (context === undefined) {
    throw new Error("useVault must be used within a VaultProvider");
  }
  return context;
}
