// import { motion, useMotionValue, useTransform } from "framer-motion";
// import { useEffect } from "react";

import { BannerButton, BannerButtonProps } from "@/primitives/Buttons";
import { cn } from "@heroui/react";
import { Image } from "./Image";
import { Disc, User } from "lucide-react";
// import { Pen, Play, Shuffle, Upload } from "lucide-react";

export default function BannerNew({
  children,
  imageUrl,
  // coverImageUrl,
  title,
  subtitle,
  bannerType,
  // scrollRef,
  buttons,
}: {
  children?: React.ReactNode;
  imageUrl: string;
  title: string;
  subtitle?: string;
  bannerType: "artist" | "album"
  // coverImageUrl?: string;
  // scrollRef: React.RefObject<HTMLElement>;
  buttons?: BannerButtonProps[];
}) {
  // const scrollY = useMotionValue(0);

  // useEffect(() => {
  //   const el = scrollRef.current;

  //   if (!el) return;
  //   const handleScroll = () => scrollY.set(el.scrollTop);

  //   el.addEventListener("scroll", handleScroll);

  //   return () => el.removeEventListener("scroll", handleScroll);
  // }, [scrollRef, scrollY]);

  // const scrollHpx = 200;

  // const opacity = useTransform(scrollY, [100, scrollHpx], [1, 0]);
  // const opacityInverse = useTransform(scrollY, [100, scrollHpx], [0, 1]);

  const image = (): React.ReactNode => {
    if (bannerType == "artist") {
      return (
          <div className="overflow-hidden relative rounded-full shrink-0 bg-default shadow-[0_40px_60px_-30px_rgba(20,40,40,0.55)] ring-[6px] ring-white/60 w-70 h-70">
            <Image
              fallbackIcon={User}
              height={70 * 0.25 * 16}
              width={70 * 0.25 * 16}
              radius="none"
              src={imageUrl}
              className="object-cover opacity-100 transition-opacity duration-300 size-full motion-reduce:transition-none"
            />
          </div>
      )
    }

    if (bannerType == "album") {
      return (
          <div className="overflow-hidden relative rounded-xl shrink-0 bg-default shadow-[0_40px_60px_-30px_rgba(20,40,40,0.55)] w-70 h-70 border-border border-1">
            <Image
              fallbackIcon={Disc}
              height={70 * 0.25 * 16}
              width={70 * 0.25 * 16}
              radius="none"
              src={imageUrl}
              className="object-cover opacity-100 transition-opacity duration-300 size-full motion-reduce:transition-none"
            />
            <img
              alt={title}
              className={cn(
                "size-full object-cover transition-opacity duration-300 motion-reduce:transition-none opacity-100",
              )}
              src={imageUrl}
            />
          </div>
      )
    }
  }

  return (
    <div className="flex-col py-5 w-full min-h-0 max-w-body">
      <div className="flex gap-2 justify-end">
            {buttons &&
              buttons.map((props: BannerButtonProps, i) => {
                return <BannerButton key={props.tooltip ?? i} {...props} />;
              })}
      </div>
      <div className="flex gap-10 mt-2">
        {image()}
        <div className="flex flex-col flex-grow gap-0 justify-end min-h-full">
          <h1 className="text-8xl font-black font-display">{title}</h1>
          <p className="text-lg text-foreground/80">{subtitle}</p>
            {children}
        </div>
      </div>
    </div>
  )

  // return (
  //   <>
  //     <motion.div
  //       className="flex overflow-hidden sticky top-0 right-0 left-0 z-50 gap-12 items-center px-4 w-full cursor-default select-none sm:px-6 group border-b-1 border-border h-[100px]"
  //       style={{
  //         opacity: opacityInverse,
  //         backgroundImage: bgImageUrl
  //           ? `linear-gradient(
  //               to bottom,
  //               color-mix(in oklch, var(--background) 80%, transparent),
  //               color-mix(in oklch, var(--background) 90%, transparent)
  //             ), url(${bgImageUrl})`
  //           : "",
  //         backgroundSize: "cover",
  //         backgroundPosition: "center",
  //       }}
  //     >
  //       <div className="flex justify-between items-center w-full">
  //         <div className="flex gap-4 items-center pl-2 w-full justify-left">
  //           {coverImageUrl && (
  //             <img
  //               className="rounded-md shadow-md max-h-[70px] border-1 border-border"
  //               src={coverImageUrl}
  //             />
  //           )}

  //           <h1 className="w-full text-xl font-extrabold sm:text-3xl text-foreground max-w-body">
  //             {title}
  //           </h1>
  //         </div>
  //         <div className="flex gap-2 justify-end items-center h-full transition sm:opacity-0 group-hover:opacity-100">
  //           {buttons &&
  //             buttons.map((props: BannerButtonProps, i) => {
  //               return <BannerButton key={props.tooltip ?? i} {...props} />;
  //             })}
  //         </div>
  //       </div>
  //     </motion.div>

  //     <motion.div
  //       className="flex overflow-hidden gap-12 items-center px-4 w-full cursor-default select-none sm:px-16 -mt-[100px] border-b-1 border-border h-[400px] group"
  //       style={{
  //         opacity,
  //         backgroundImage: bgImageUrl
  //           ? `linear-gradient(
  //               to bottom,
  //               color-mix(in oklch, var(--background) 80%, transparent),
  //               color-mix(in oklch, var(--background) 90%, transparent)
  //             ), url(${bgImageUrl})`
  //           : "",
  //         backgroundSize: "cover",
  //         backgroundPosition: "center",
  //       }}
  //     >
  //       <div className="flex w-full h-full justify-left">
  //         <div className="flex flex-col gap-0 justify-start items-center pt-4 w-full sm:flex-row sm:gap-10 sm:pt-0 max-w-body">
  //           {coverImageUrl && (
  //             <img
  //               className="rounded-lg shadow-md border-1 border-border max-h-[250px] sm:max-h-[280px]"
  //               src={coverImageUrl}
  //             />
  //           )}

  //           <div className="flex relative flex-col w-full h-full sm:static">
  //             <div className="pt-4 sm:pt-0 sm:h-1/4" />
  //             <div className="h-4/5 sm:h-1/2">
  //               <h1 className="w-full text-xl font-extrabold lg:text-3xl xl:text-6xl text-foreground max-w-body">
  //                 {title}
  //               </h1>
  //               <div className="flex overflow-hidden z-10 flex-col h-full sm:gap-3 sm:pt-3 text-foreground max-w-body">
  //                 {children}
  //               </div>
  //             </div>
  //             <div className="flex absolute right-0 bottom-0 gap-2 justify-end items-end pb-8 h-1/4 transition sm:static sm:opacity-0 group-hover:opacity-100">
  //               {buttons &&
  //                 buttons.map((props: BannerButtonProps, i) => {
  //                   return <BannerButton key={props.tooltip ?? i} {...props} />;
  //                 })}
  //             </div>
  //           </div>
  //         </div>
  //       </div>
  //     </motion.div>
  //   </>
  // );
}
