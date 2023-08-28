import { Handle, Position } from "reactflow";
import { IUtils } from "./interfaces";

interface IProps {
  data: { label: string; cpu: IUtils; memory: IUtils; ins: string[] };
}

const ServiceNode = ({ data }: IProps) => {
  return (
    <>
      {data.ins.map((h, i) => (
        <Handle
          key={i}
          type="target"
          id={h}
          position={Position.Left}
          style={{
            top: `${(i + 1) * (100 / (data.ins.length + 1))}%`,
          }}
        />
      ))}

      <div
        style={{
          backgroundColor: `hsl(${100 - data.cpu.quantile * 100}, 100%, 50%)`,
          border: "black solid 1px",
          padding: "0.5em",
          borderRadius: "0.5em",
        }}
      >
        <p style={{ textAlign: "center", margin: "0" }}>{data.label}</p>
        <table>
          <tr>
            <th></th>
            <th>Mean</th>
            <th>Quantile</th>
            <th>Stdev</th>
          </tr>
          <tr>
            <th>CPU</th>
            <td>{(data.cpu.mean * 100).toFixed(2)}%</td>
            <td>{(data.cpu.quantile * 100).toFixed(2)}%</td>
            <td>{data.cpu.stdev.toFixed(4)}</td>
          </tr>
          <tr>
            <th>Memory</th>
            <td>{(data.memory.mean * 100).toFixed(2)}%</td>
            <td>{(data.memory.quantile * 100).toFixed(2)}%</td>
            <td>{data.memory.stdev.toFixed(4)}</td>
          </tr>
        </table>
      </div>
      <Handle type="source" position={Position.Right} id="a" />
    </>
  );
};

export default ServiceNode;
