import React, { useCallback, useMemo } from "react";
import Dagre from "@dagrejs/dagre";
import ReactFlow, {
    Edge,
    Node,
    Panel,
    useEdgesState,
    useNodesState,
    useReactFlow,
} from "reactflow";
import { flowChart } from "./utils/getYChart";

import "reactflow/dist/style.css";
import ServiceNode from "./ServiceNode";

const g = new Dagre.graphlib.Graph().setDefaultEdgeLabel(() => ({}));

const getLayoutedElements = (nodes: Node[], edges: Edge[], options: any) => {
    g.setGraph({ rankdir: options.direction });

    edges.forEach((edge) => g.setEdge(edge.source, edge.target));
    nodes.forEach((node) => g.setNode(node.id, node));

    Dagre.layout(g);

    return {
        nodes: nodes.map((node) => {
            const { x, y } = g.node(node.id);

            return { ...node, position: { x, y } };
        }),
        edges,
    };
};

const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
    flowChart.nodes,
    Object.values(flowChart.edges).flat(),
    { direction: "TB" }
);

const LayoutFlow = () => {
    const nodeTypes = useMemo(() => ({ serviceNode: ServiceNode }), []);

    const { fitView } = useReactFlow();
    const [nodes, setNodes, onNodesChange] = useNodesState([...layoutedNodes]);
    const [edges, setEdges, onEdgesChange] = useEdgesState([...layoutedEdges]);

    const onLayout = useCallback(
        (direction: string) => {
            const layouted = getLayoutedElements(nodes, edges, { direction });

            setNodes([...layouted.nodes]);
            setEdges([...layouted.edges]);

            window.requestAnimationFrame(() => {
                fitView();
            });
        },
        [nodes, edges]
    );

    const handleChangeEdges = (key: string) => {
        if (key === "") return setEdges(Object.values(flowChart.edges).flat());
        setEdges(flowChart.edges[key]);
    };

    return (
        <ReactFlow
            nodeTypes={nodeTypes}
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            fitView
            style={{ flexGrow: 1 }}
        >
            <Panel
                position="top-right"
                style={{
                    display: "flex",
                    flexDirection: "row",
                    gap: "0.5em",
                    backgroundColor: "#475569",
                    padding: "0.5em",
                    borderRadius: "0.5em",
                }}
            >
                <button onClick={() => onLayout("TB")}>vertical layout</button>
                <button onClick={() => onLayout("LR")}>
                    horizontal layout
                </button>
            </Panel>
            <Panel
                position="bottom-right"
                style={{
                    display: "flex",
                    flexDirection: "row",
                    gap: "0.5em",
                    backgroundColor: "#475569",
                    padding: "0.5em",
                    borderRadius: "0.5em",
                }}
            >
                <button onClick={() => handleChangeEdges("")}>Show All</button>
                {Object.keys(flowChart.edges).map((e, i) => (
                    <button key={i} onClick={() => handleChangeEdges(e)}>
                        Show {e}
                    </button>
                ))}
            </Panel>
        </ReactFlow>
    );
};

export default LayoutFlow;
