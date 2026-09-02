export enum Event {
  MsgErr = "msg.error",
  MsgSuccess = "msg.success",
  MsgInfo = "msg.info",
  MsgWarn = "msg.warn",

  ServerMessage = "server.message",

  AuthLoggedIn = "auth.loggedin",

  ServerOnline = "server.online",
  ServerStatus = "server.status",
  SyncInProgress = "sync.inprogress",

  PlayerContextChange = "player.contextchange",
  PlayerPlaybackChange = "player.playbackchange",

  PlayerLocalContextChange = "player.local.contextchange",
  PlayerLocalPlaybackChange = "player.local.playbackchange",

  PlayerLocalStartContext = "player.local.startcontext",
  PlayerLocalPlayPause = "player.local.playpause",
  PlayerLocalNextTrack = "player.local.nexttrack",

  DeviceConnectionID = "device.connection.id",

  DeviceConnected = "device.connected",
  DeviceUpdated = "device.updated",
  DeviceDisconnected = "device.disconnected",
  DevicePlayContext = "device.playcontext",
  DevicePlaybackState = "device.playbackstate",

  EntitiesUpdated = "entities.updated",

  PlayerPlayContext = "player.playcontext",
  PlayerPlaybackState = "player.playbackstate",
  PlayerStart = "player.start",
  PlayerPlayPause = "player.playpause",

  UserUpdated = "user.updated",

  InternalClientAllReplaced = "internal.client.allreplaced",

  ImporterItemsUpdated = "importer.itemsupdated"
}

export type Message = {
  title: string;
  message: string;
};

export type ServerMessage = {
  level: "info"
  title: string;
  body: string;
};

export type SSEvent = {
  id: string;
  type: Event;
  data?: object;
}
