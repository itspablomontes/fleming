import { createConfig, http } from "wagmi";
import { localhost, mainnet } from "wagmi/chains";
import { injected } from "wagmi/connectors";

export const config = createConfig({
	chains: [mainnet, localhost],
	connectors: [injected()],
	transports: {
		[mainnet.id]: http(),
		[localhost.id]: http(),
	},
});
