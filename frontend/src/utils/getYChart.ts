import { Edge, MarkerType, Node, Position } from "reactflow";
import { IYChart } from "../interfaces";
import data from "../y-chart.json";
import Dagre from "@dagrejs/dagre";

export const getYChart = () => {
    return data as IYChart;
};

export const getFlowChart = () => {
    const data = getYChart();
    const services: Node[] = [];
    const operations: { [key: string]: Edge[] } = {};

    Object.keys(data.services).map((k, i) => {
        console.log(k);
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

    return { nodes: services, edges: operations };
};

export const flowChart = getFlowChart();
