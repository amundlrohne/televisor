import { ReactFlowProvider } from "reactflow";

import "reactflow/dist/style.css";
import LayoutFlow from "./LayoutFlow";
import { getYChart } from "./utils/getYChart";

export default function App() {
    const data = getYChart();

    return (
        <ReactFlowProvider>
            <div
                style={{
                    width: "100vw",
                    height: "100vh",
                    display: "flex",
                    flexDirection: "row",
                    backgroundColor: "#fefefe",
                }}
            >
                <div
                    style={{
                        backgroundColor: "#475569",
                        display: "flex",
                        flexDirection: "column",

                        overflowY: "auto",

                        gap: "1em",
                        padding: "1em",
                        borderRadius: "0 1em 1em 0",
                    }}
                >
                    {data.annotations.map((a, i) => (
                        <div
                            style={{
                                backgroundColor: "#fefefe",
                                padding: "0 1em 0 1em",
                            }}
                            key={i}
                        >
                            <h4>{a.annotationType}</h4>
                            <p>{a.message}</p>
                            <p>Services:</p>
                            <ul>
                                {a.services?.map((s, si) => (
                                    <li key={"s" + si}>{s}</li>
                                ))}
                            </ul>
                            <p>Operations:</p>
                            <ul>
                                {a.operations?.map((o, oi) => (
                                    <li key={"o" + oi}>{o}</li>
                                ))}
                            </ul>
                        </div>
                    ))}
                </div>
                <LayoutFlow />
            </div>
        </ReactFlowProvider>
    );
}
