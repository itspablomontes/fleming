import { Activity, ArrowRight, Fingerprint, Microscope } from "lucide-react";
import { motion } from "motion/react";
import { cn } from "@/lib/utils";

/**
 * UseCaseCard component
 * Displays a single use case with icon, title, description and CTA
 */
interface UseCaseCardProps {
    icon: React.ElementType;
    title: string;
    description: string;
    cta: string;
    delay: number;
    color: string;
}

function UseCaseCard({ icon: Icon, title, description, cta, delay, color }: UseCaseCardProps) {
    return (
        <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, margin: "-50px" }}
            transition={{ duration: 0.5, delay }}
            className="group relative overflow-hidden rounded-3xl border bg-background/50 p-8 hover:bg-background/80 transition-colors"
        >
             {/* Gradient blob for hover effect */}
             <div className={cn(
                "absolute -right-20 -top-20 h-64 w-64 rounded-full bg-linear-to-br opacity-0 blur-3xl transition-opacity duration-500 group-hover:opacity-10",
                color
            )} />
            
            <div className="relative z-10 flex flex-col h-full">
                <div className={cn(
                    "mb-6 flex h-12 w-12 items-center justify-center rounded-2xl bg-muted transition-colors group-hover:bg-primary/10",
                    "text-muted-foreground group-hover:text-primary"
                )}>
                    <Icon className="h-6 w-6" />
                </div>
                
                <h3 className="mb-3 text-xl font-bold tracking-tight">{title}</h3>
                
                <p className="mb-8 text-muted-foreground leading-relaxed grow">
                    {description}
                </p>
                
                <div className="flex items-center text-sm font-medium text-primary cursor-pointer group/link">
                    {cta}
                    <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover/link:translate-x-1" />
                </div>
            </div>
        </motion.div>
    );
}

/**
 * UseCasesSection
 * Visualizes practical applications: Proving (ZK), Research (DeSci), Tracking (Protocols).
 */
export function UseCasesSection() {
    const cases = [
        {
            icon: Fingerprint,
            title: "Prove Your Health",
            description: "Generate Zero-Knowledge proofs of your biological age without revealing your birthday. Verify you meet trial criteria without exposing your medical history.",
            cta: "Explore ZK Proofs",
            color: "from-blue-500 to-cyan-500"
        },
        {
            icon: Microscope,
            title: "Contribute to Science",
            description: "Join decentralized clinical trials (DeSci) with a single click. Your data, your terms.",
            cta: "Join the Network",
            color: "from-emerald-500 to-green-500"
        },
        {
            icon: Activity,
            title: "Track Protocols",
            description: "Log your longevity interventions—supplements, fasts, exercises—and visualize their impact on your biomarkers in real-time. Own your N=1 experiment.",
            cta: "Start Tracking",
            color: "from-purple-500 to-pink-500"
        }
    ];

    return (
        <section className="py-16 md:py-24 bg-muted/20 relative">
            <div className="container px-4 mx-auto max-w-6xl">
                 <div className="text-center mb-10 md:mb-12 space-y-4">
                    <h2 className={cn(
                        "text-3xl md:text-4xl font-bold tracking-tight inline-block",
                        "bg-clip-text text-transparent bg-linear-to-r from-foreground to-foreground/70"
                    )}>
                        Utility First
                    </h2>
                    <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
                        Fleming isn't just a vault. It's a tool for action.
                    </p>
                </div>

                <div className="grid md:grid-cols-3 gap-6 lg:gap-8">
                    {cases.map((useCase, index) => (
                        <UseCaseCard 
                            key={useCase.title}
                            {...useCase}
                            delay={index * 0.1}
                        />
                    ))}
                </div>
            </div>
        </section>
    );
}
