import { useSession } from "@/session/SessionContext";
import Banner from "@/components/Banner";

export default function HomePage({
  scrollRef,
}: {
  scrollRef: React.RefObject<HTMLElement>;
}) {
  const session = useSession();

  const getGreeting = (): string => {
    const hour = new Date().getHours();
    if (hour < 12) {
      return "Good Morning";
    }
    if (hour < 18) {
      return "Good Afternoon";
    }
    return "Good Evening";
  };

  return (
    <>
      <Banner
        bgImageUrl="https://upload.wikimedia.org/wikipedia/commons/9/9b/Early_Summer_Morning_%286786835628%29.jpg"
        scrollRef={scrollRef}
        title={`${getGreeting()} ${session.user?.username}`}
      >
        <p></p>
      </Banner>
      <div className="flex flex-col pt-4">
      {
        // <p className="pt-8 text-lg">Yesterday something was annoying you</p>
        // <div className="flex gap-4 px-4 pt-2">
        //   <MoodCard mood="angry" frac={0.72}/>
        //   <MoodCard mood="sad" frac={0.18}/>
        //   <MoodCard mood="calm" frac={0.08}/>
        //   <MoodCard mood="happy" frac={0.02}/>
        // </div>

        // <p className="pt-8 text-lg">But last month has been good for you!</p>
        // <div className="flex gap-4 px-4 pt-2">
        //   <MoodCard mood="happy" frac={0.51}/>
        //   <MoodCard mood="calm" frac={0.28}/>
        //   <MoodCard mood="angry" frac={0.17}/>
        //   <MoodCard mood="sad" frac={0.04}/>
        // </div>

        // <p className="pt-8 text-lg">Maybe you need to listen to some happy music today to make you feel better?</p>
        }
      </div>
      </>
  );
}
