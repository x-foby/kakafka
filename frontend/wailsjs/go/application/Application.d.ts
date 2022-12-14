// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {kafka} from '../models';
import {application} from '../models';

export function Connect(arg1:string):Promise<Error>;

export function ConsumerOffsets(arg1:string,arg2:string):Promise<Array<kafka.ConsumerOffset>>;

export function CreateProfile(arg1:application.Profile):Promise<Error>;

export function CreateTopic(arg1:string,arg2:kafka.TopicConfig):Promise<kafka.Topic>;

export function DeleteProfile(arg1:string):Promise<Error>;

export function DeleteTopic(arg1:string,arg2:string):Promise<Error>;

export function GetConfigs():Promise<application.Config>;

export function GetTopics(arg1:string,arg2:boolean):Promise<Array<kafka.Topic>>;
