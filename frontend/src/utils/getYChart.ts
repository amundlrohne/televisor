import { Edge, MarkerType, Node, Position } from "reactflow";
import { IYChart } from "../interfaces";
import detection from "../y-chart.json";
import recommendation from "../y-chart-recommendation.json";

export const getDetectionSDG = () => {
  return detection as IYChart;
};

export const getRecommendationSDG = () => {
  return recommendation as IYChart;
};

export const getFlowChart = (ychart: IYChart) => {
  const data = ychart;
  const services: Node[] = [];
  const operations: { [key: string]: Edge[] } = {};

  Object.keys(data.services).map((k, i) => {
    return services.push({
      id: k,
      position: { x: 0, y: i * 100 },
      data: {
        label: k,
        cpu: data.services[k].cpu,
        memory: data.services[k].memory,
        ins: [],
      },
      sourcePosition: Position.Right,
      targetPosition: Position.Left,
      type: "serviceNode",
    });
  });

  Object.keys(data.operations).map((k, i) => {
    const operation: Edge[] = [];
    Object.keys(data.operations[k]).map((kk, ii) => {
      const edge = data.operations[k][kk];
      if (edge.to !== edge.from) {
        const si = services.findIndex((s, i) => s.id === edge.to);
        const inputId = `${k}-${kk}`;
        services[si].data.ins.push(inputId);

        return operation.push({
          id: `${kk}-${edge.from}-${edge.to}`,
          source: edge.from,
          target: edge.to,
          targetHandle: inputId,
          markerEnd: {
            type: MarkerType.Arrow,
          },
          label: kk,
        });
      }
    });
    operations[k] = operation;
  });

  services.forEach((s) => console.log(s));

  return { nodes: services, edges: operations };
};

export const detectionSDG = getFlowChart(getDetectionSDG());
export const recommendationSDG = getFlowChart(getRecommendationSDG());
