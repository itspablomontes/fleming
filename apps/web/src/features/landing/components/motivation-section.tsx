import { motion } from "motion/react";

/**
 * MotivationSection
 * A minimalist, philosophical section centered on the concepts of health sovereignty and leverage.
 * Mimics the writing style of Naval Ravikant (aphoristic, punchy, high-leverage)
 * but gathered from VitaDAO ethos (community ownership, open science, aging as a modifiable process).
 */
export function MotivationSection() {
	return (
		<section className="py-24 md:py-32 bg-background relative overflow-hidden">
			{/* Subtle background gradient for depth - simplified, minimalist */}
			<div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,var(--tw-gradient-stops))] from-primary/5 via-background to-background opacity-40 pointer-events-none" />

			<div className="container px-4 mx-auto max-w-6xl relative z-10">
				{/* Section Header */}
				<div className="text-center mb-16">
					<motion.h2
						initial={{ opacity: 0, y: 10 }}
						whileInView={{ opacity: 1, y: 0 }}
						viewport={{ once: true }}
						transition={{ duration: 0.6 }}
						className="text-xs font-mono uppercase tracking-[0.3em] text-muted-foreground mb-4"
					>
						The Philosophy
					</motion.h2>
					<motion.h3
						initial={{ opacity: 0, scale: 0.95 }}
						whileInView={{ opacity: 1, scale: 1 }}
						viewport={{ once: true }}
						transition={{ duration: 0.8, ease: "easeOut" }}
						className="text-3xl md:text-5xl font-serif font-medium tracking-tight leading-tight bg-clip-text text-transparent bg-linear-to-b from-foreground to-foreground/70"
					>
						The Long Game
					</motion.h3>
				</div>

				{/* Single Paragraph Content */}
				<motion.div
					initial={{ opacity: 0, y: 20 }}
					whileInView={{ opacity: 1, y: 0 }}
					viewport={{ once: true, margin: "-50px" }}
					transition={{ duration: 0.8, delay: 0.2, ease: "easeOut" }}
					className="relative"
				>
					{/* Decorative accent lines */}
					<div className="absolute -top-8 left-1/2 -translate-x-1/2 w-24 h-px bg-linear-to-r from-transparent via-primary/30 to-transparent" />
					<div className="absolute -bottom-8 left-1/2 -translate-x-1/2 w-24 h-px bg-linear-to-r from-transparent via-primary/30 to-transparent" />

					<p className="text-lg md:text-xl lg:text-2xl font-normal leading-relaxed text-center text-foreground/90 max-w-5xl mx-auto">
						Your biological data is the map to your own longevity, yet it remains
						fragmented, siloed, and sold by intermediaries. Fleming restores the map
						to its explorer. We replace institutional trust with{" "}
						<span className="text-foreground font-medium border-b-2 border-primary/20">
							cryptographic truth
						</span>
						, building a sovereign vault where your medical history is secured by
						mathematics, not policy. This enables a paradox: you can now prove your
						health to the world without ever revealing your identity. By contributing
						to open science on your own terms, you don't just protect your
						privacy, you accelerate the cure. To win against entropy, we must first
						own the data.
					</p>
				</motion.div>
			</div>
		</section>
	);
}
