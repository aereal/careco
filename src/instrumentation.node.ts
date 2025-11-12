import { propagation } from '@opentelemetry/api';
import { W3CTraceContextPropagator } from '@opentelemetry/core';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { resourceFromAttributes } from '@opentelemetry/resources';
import {
  NodeTracerProvider,
  SimpleSpanProcessor,
} from '@opentelemetry/sdk-trace-node';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';
import { FetchInstrumentation } from '@vercel/otel';

export const register = async (): Promise<void> => {
  const exporter = new OTLPTraceExporter();
  const tp = new NodeTracerProvider({
    resource: resourceFromAttributes({ [ATTR_SERVICE_NAME]: 'careco-front' }),
    spanProcessors: [new SimpleSpanProcessor(exporter)],
  });
  const propagator = new W3CTraceContextPropagator();
  tp.register({ propagator });
  propagation.setGlobalPropagator(propagator);

  registerInstrumentations({
    instrumentations: [new FetchInstrumentation()],
  });
};
