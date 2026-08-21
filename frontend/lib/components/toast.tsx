import { useEffect } from "react";
import { toast, ToastProvider } from "@heroui/react";

import { Event, Message, ServerMessage } from "@/types/events";
import { useApi } from "@/api/ApiContext";

export default function Toast() {
  const api = useApi();

  useEffect(() => {
    const serverMsgHandler = (msg: ServerMessage) => {
      let variant: "success" | "danger" = "danger";

      switch (msg.level) {
        case "info":
          variant = "success";
          break;
      }

      toast(msg.title, {
        description: msg.body,
        variant,
      });
    };

    const errHandler = (msg: Message) => {
      console.log(msg);

      toast(msg.title, {
        description: msg.message,
        variant: "danger",
      });
    };

    const successHandler = (e: string) => {
      toast("Success", {
        description: e,
        variant: "success",
      });
    };

    const infoHandler = (e: string) => {
      toast("Information", {
        description: e,
        variant: "accent",
      });
    };

    const warnHandler = (e: string) => {
      toast("Warning", {
        description: e,
        variant: "warning",
      });
    };

    const unsub1 = api.eventSubscribe(Event.MsgErr, errHandler);
    const unsub2 = api.eventSubscribe(Event.MsgSuccess, successHandler);
    const unsub3 = api.eventSubscribe(Event.MsgInfo, infoHandler);
    const unsub4 = api.eventSubscribe(Event.MsgWarn, warnHandler);
    const unsub5 = api.eventSubscribe(Event.ServerMessage, serverMsgHandler);

    return () => {
      unsub1()
      unsub2()
      unsub3()
      unsub4()
      unsub5()
    };
  }, [api]);

  return <ToastProvider placement="top" />;
}
