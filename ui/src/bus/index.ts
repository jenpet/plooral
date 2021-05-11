// To resolve the dependency hell this index file just exports core elements and the initialization
import {init} from "@/bus/service";
import {getBus, BusEvent, ComponentBus } from "@/bus/core";

export {init as initBus, getBus, BusEvent, ComponentBus}