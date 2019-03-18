import {Injectable} from '@angular/core';
import {CELoadedEvent, EventType} from "../../models/events";
import {fromEvent, Observable} from "rxjs";
import {publish} from "rxjs/operators";

@Injectable()
export class CeRegisterService {
  // <componentEntryPoint, componentRef>
  readonly registeredCEs: Map<string, string>;
  readonly bus: Element;
  readonly ceLoadedEvents$: Observable<CELoadedEvent>;

  constructor() {
    this.registeredCEs = new Map<string, string>();
    this.bus = document.querySelector<Element>('head');

    // this.ceLoadedEvents$ = fromEvent<CELoadedEvent>(this.bus, EventType.CE_LOADED_EVENT);

    new Observable((observer) => {
      // Get the next and error callbacks. These will be passed in when
      // the consumer subscribes.
      const {next, error} = observer;
      let watchId;

      // Simple geolocation API check provides values to publish
      if ('geolocation' in navigator) {
        watchId = navigator.geolocation.watchPosition(next, error);
      } else {
        error('Geolocation not available');
      }

      // When the consumer unsubscribes, clean up data ready for next subscription.
      return {unsubscribe() { navigator.geolocation.clearWatch(watchId); }};
    });

    this.ceLoadedEvents$ = new Observable((observer) => {
      const {next, error} = observer;
      let newEvent: CELoadedEvent;

      let handler = (msg: CustomEvent<CELoadedEvent>) => {
        //get event that new custom element was loaded
        console.log(msg);
        next(msg.detail);
      };

      this.bus.addEventListener(EventType.CE_LOADED_EVENT, handler);
      console.log(EventType.CE_LOADED_EVENT);
      return {unsubscribe() { this.bus.removeEventListener(EventType.CE_LOADED_EVENT, handler) }};

    });

  }

  public registerCe(componentEntryPoint: string): Observable<CELoadedEvent> {

    const script = document.createElement('script');
    script.src = 'http://54.183.122.86:32291' + componentEntryPoint;
    console.log(script);

    this.ceLoadedEvents$.subscribe((event: CELoadedEvent) => {
        console.debug(event);
    });

    // debug(getEventListeners(document.querySelector('head')).CELoadedEvent[0].listener);

    this.bus.appendChild(script);



    return this.ceLoadedEvents$;

  }



}
