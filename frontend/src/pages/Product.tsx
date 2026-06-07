import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
export interface Product {
    ProductID: number;
    Name: string;
    image_url: string;
    description: string;
    SourceURL: string;
    LastPrice: number;
    LowestPrice: number;
    CreatedBy: number;
    CreateAt: string;
    UpdatedAt: string;
}
export interface PriceHistory {
    Price: number;
    RegisterDate: string;
}

export default function ProductView() {
    const { id } = useParams();

    const [product, setProduct] = useState<Product | null>(null);
    const [loading, setLoading] = useState(true);
    const [targetPrice, setTargetPrice] = useState("");
    const [isFollowing, setIsFollowing] = useState(false);
    const [savedTargetPrice, setSavedTargetPrice] = useState<number | null>(null);
    const [isUpdating, setIsUpdating] = useState(false);
    const [progress, setProgress] = useState(0);
    const [, setMessage] = useState("");
    const [, setOpacity] = useState(0);
    const [history, setHistory] = useState<PriceHistory[]>([]);
    const [timeRange, setTimeRange] = useState(7);

    const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

    const showMessage = (text: string) => {
        setMessage(text);
        setOpacity(1);
        setTimeout(() => setOpacity(0), 3000);
        setTimeout(() => setMessage(""), 3500);
    };

    useEffect(() => {
        const loadProductData = async () => {
            try {
                setLoading(true);
                setIsFollowing(false);
                setSavedTargetPrice(null);

                const response = await fetch(`${API_URL}/api/products/${id}`, { credentials: "include" });
                if (!response.ok) throw new Error("Producto no encontrado");

                const productData = await response.json();
                setProduct(productData);

                const trackingResponse = await fetch(`${API_URL}/api/check-tracking/${productData.ProductID}`, {
                    credentials: "include"
                });

                if (trackingResponse.ok) {
                    const trackData = await trackingResponse.json();
                    setIsFollowing(true);
                    setIsFollowing(!!trackData.is_following);

                    if (trackData.target_price && trackData.target_price > 0) {
                        setSavedTargetPrice(trackData.target_price);
                        setTargetPrice(trackData.target_price.toString());
                    } else {
                        setTargetPrice(productData.LowestPrice?.toString() || "");
                    }
                }
            } catch (error) {
                console.error("Error al inicializar vista:", error);
            } finally {
                setLoading(false);
            }
        };

        if (id) loadProductData();
    }, [id]);

    const fetchHistory = async (days: number) => {
        try {
            const res = await fetch(`${API_URL}/api/products/${id}/history?days=${days}`, {
                credentials: "include"
            });

            if (!res.ok) throw new Error("No autorizado");

            const data = await res.json();
            setHistory(Array.isArray(data) ? data : []);
        } catch (error) {
            console.error("Error al cargar historial:", error);
            setHistory([]);
        }
    };
    useEffect(() => {
        fetchHistory(timeRange);
    }, [id, timeRange]);


    const handleFollowProduct = async () => {
        if (!product) return;

        try {
            if (isFollowing) {
                const response = await fetch(`${API_URL}/api/tracking/${product.ProductID}`, {
                    method: "DELETE",
                    credentials: "include"
                });

                if (response.ok) {
                    setIsFollowing(false);
                    setSavedTargetPrice(null);
                } else {
                    console.error("Error al borrar el seguimiento:", response.status);
                }
            } else {
                const response = await fetch(`${API_URL}/api/tracking`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    credentials: "include",
                    body: JSON.stringify({
                        product_id: product.ProductID,
                        target_price: 0,
                        notify_price_changes: true
                    })
                });
                if (response.ok) {
                    setIsFollowing(true);
                }
            }
        } catch (error) {
            console.error("Error al alternar seguimiento:", error);
        }
    };

    const handleUpdatePrice = async () => {
        if (!product) return;
        setIsUpdating(true);
        setProgress(30);

        try {
            const response = await fetch(`${API_URL}/api/track`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ url: product.SourceURL })
            });

            if (response.ok) {
                setProgress(100);
                const updatedData = await response.json();
                setProduct(updatedData.data);
            }
        } catch (error) {
            console.error("Error:", error);
        } finally {
            setTimeout(() => {
                setIsUpdating(false);
                setProgress(0);
            }, 500);
        }
    };

    const handleCreateAlert = async () => {
        if (!product) return;

        const priceToSet = targetPrice ? parseFloat(targetPrice) : product.LowestPrice;

        try {
            const response = await fetch(`${API_URL}/api/tracking`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({
                    product_id: product.ProductID,
                    target_price: priceToSet,
                    notify_price_changes: true
                })
            });

            if (response.ok) {
                setSavedTargetPrice(priceToSet);
                setIsFollowing(true);
                showMessage(`Alerta establecida en ${priceToSet}€`);
            }
        } catch {
            showMessage("Error al guardar la alerta");
        }
    };

    if (loading) {
        return (
            <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh", backgroundColor: "#fafafa" }}>
                <Header />
                <main style={{ flex: 1, display: "flex", justifyContent: "center", alignItems: "center" }}>
                    <h2>Cargando producto... ⏳</h2>
                </main>
                <Footer />
            </div>
        );
    }

    if (!product) return <h2>Producto no encontrado</h2>;

    return (
        <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh", backgroundColor: "#fafafa" }}>
            <Header />
            <main style={{ flex: 1, display: "flex", justifyContent: "center", padding: "2rem" }}>
                <div style={{ maxWidth: "900px", width: "100%", display: "flex", flexDirection: "column", gap: "2rem", backgroundColor: "#fff", padding: "2rem", borderRadius: "8px", boxShadow: "0 4px 6px rgba(0,0,0,0.05)" }}>

                    <div style={{ display: "flex", gap: "2rem", alignItems: "flex-start", flexWrap: "wrap" }}>
                        <div style={{ flex: "1 1 30%", minWidth: "250px" }}>
                            <img src={product.image_url} alt={product.Name} style={{ width: "100%", borderRadius: "8px", objectFit: "contain", border: "1px solid #f0f0f0" }} />
                        </div>
                        <div style={{ flex: "2 1 60%" }}>
                            <h1 style={{ margin: "0 0 1rem 0", fontSize: "2.2rem", color: "#111827", lineHeight: "1.2" }}>{product.Name}</h1>
                        </div>
                    </div>

                    {/* Action Buttons */}
                    <div style={{ display: "flex", gap: "1rem" }}>
                        <button
                            onClick={handleFollowProduct}
                            style={{
                                ...btnStylePrimary,
                                backgroundColor: isFollowing ? "#111827" : "#FACC15",
                                color: isFollowing ? "#FACC15" : "#111827",
                                border: isFollowing ? "2px solid #FACC15" : "none"
                            }}
                            title={isFollowing ? "Dejar de seguir" : "Seguir producto"}
                        >
                            <svg width="24" height="24" viewBox="0 0 24 24" fill={isFollowing ? "currentColor" : "none"} stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z" />
                            </svg>
                        </button>
                        <button
                            className="btn-wrapper"
                            style={{ ...btnStyleSecondary, position: "relative", overflow: "hidden" }}
                            onClick={handleUpdatePrice}
                            disabled={isUpdating}
                            title="Actualizar precio"
                        >
                            {isUpdating && <div className="progress-bar-full" style={{ width: `${progress}%` }}></div>}
                            <span className={isUpdating ? "spinning" : ""} style={{ display: "inline-block" }}>
                                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                                    <path d="M23 4v6h-6M1 20v-6h6M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15" />
                                </svg>
                            </span>
                        </button>
                        <a
                            href={product.SourceURL}
                            target="_blank"
                            rel="noopener noreferrer"
                            style={btnStyleLink}
                            title="Ir al sitio web"
                        >
                            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6M15 3h6v6M10 14L21 3" />
                            </svg>
                        </a>
                    </div>

                    <hr style={{ border: "none", borderTop: "1px solid #eaeaea", margin: "0" }} />

                    {/* Price & Alerts */}
                    <div style={{ display: "flex", gap: "2rem", flexWrap: "wrap" }}>
                        <div style={{ flex: "1 1 45%", display: "flex", flexDirection: "column", gap: "0.75rem", justifyContent: "center" }}>
                            <div style={{ fontSize: "1.75rem", fontWeight: "bold", color: "#111827" }}>{product.LastPrice} €</div>
                            <div style={{ fontSize: "1rem", color: "#059669", fontWeight: "600" }}>Histórico más bajo: {product.LowestPrice} €</div>
                        </div>

                        <div style={{ flex: "1 1 45%", display: "flex", flexDirection: "column", gap: "0rem", justifyContent: "center" }}>
                            <label style={{ fontWeight: "600", color: "#4B5563", fontSize: "0.95rem", marginBottom: "0.5rem" }}>
                                Establecer precio de aviso (€):
                            </label>

                            <div style={{ display: "flex", gap: "0.5rem", alignItems: "flex-start" }}>

                                <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
                                    <input
                                        type="number"
                                        value={targetPrice}
                                        onChange={(e) => setTargetPrice(e.target.value)}
                                        placeholder={`${product.LowestPrice}`}
                                        style={{
                                            padding: "0.75rem",
                                            borderRadius: savedTargetPrice ? "6px 6px 0 0" : "6px",
                                            border: "1px solid #D1D5DB",
                                            width: "100%",
                                            boxSizing: "border-box"
                                        }}
                                    />

                                    {savedTargetPrice && (
                                        <div className="alert-active-block" style={{ width: "100%", boxSizing: "border-box" }}>
                                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
                                                <path d="M13.73 21a2 2 0 0 1-3.46 0" />
                                            </svg>
                                            Alerta configurada
                                        </div>
                                    )}
                                </div>

                                <button style={btnStyleAlert} onClick={handleCreateAlert}>
                                    {savedTargetPrice ? "Actualizar" : "Crear"}
                                </button>
                            </div>
                        </div>
                    </div>

                    <hr style={{ border: "none", borderTop: "1px solid #eaeaea", margin: "0" }} />

                    {/* Description */}
                    <div>
                        <h3>Descripción del producto</h3>
                        <p style={{ whiteSpace: "pre-wrap", color: "#4B5563", lineHeight: "1.7", backgroundColor: "#F9FAFB", padding: "1.5rem", borderRadius: "8px" }}>{product.description}</p>
                    </div>

                    {/* Graph */}
                    <div style={{ marginTop: "2rem", padding: "1.5rem", backgroundColor: "#fff", borderRadius: "8px" }}>
                        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "1rem" }}>
                            <h3>Historial de precios</h3>
                            <div>
                                <button
                                    onClick={() => setTimeRange(7)}
                                    style={{ fontWeight: timeRange === 7 ? "bold" : "normal", marginRight: "10px" }}
                                >7 Días</button>
                                <button
                                    onClick={() => setTimeRange(30)}
                                    style={{ fontWeight: timeRange === 30 ? "bold" : "normal" }}
                                >30 Días</button>
                            </div>
                        </div>
                        <div style={{ width: "100%", height: "300px" }}>
                            {loading ? (
                                <p>Cargando datos...</p>
                            ) : history.length > 0 ? (
                                <ResponsiveContainer width="100%" height="100%">
                                    <LineChart data={history}>
                                        <XAxis
                                            dataKey="RegisterDate"
                                            tickFormatter={(str) => (str && str.length >= 10 ? str.slice(5, 10) : str)}
                                        />
                                        <YAxis domain={['auto', 'auto']} />
                                        <Tooltip />
                                        <Line type="monotone" dataKey="Price" stroke="#FACC15" strokeWidth={3} />
                                    </LineChart>
                                </ResponsiveContainer>
                            ) : (
                                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%' }}>
                                    <p>No hay historial disponible para este periodo.</p>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </main>
            <Footer />
        </div>
    );
}

const btnStylePrimary = { display: "flex", alignItems: "center", justifyContent: "center", gap: "0.5rem", padding: "0.85rem 1.5rem", borderRadius: "6px", cursor: "pointer", fontWeight: "bold", fontSize: "1rem", flex: 1, transition: "0.2s" };
const btnStyleSecondary = { display: "flex", alignItems: "center", justifyContent: "center", gap: "0.5rem", padding: "0.85rem 1.5rem", backgroundColor: "#F3F4F6", color: "#374151", border: "1px solid #E5E7EB", borderRadius: "6px", cursor: "pointer", fontWeight: "bold", fontSize: "1rem", flex: 1 };
const btnStyleLink = { display: "flex", alignItems: "center", justifyContent: "center", gap: "0.5rem", padding: "0.85rem 1.5rem", backgroundColor: "#F3F4F6", color: "#374151", border: "1px solid #E5E7EB", borderRadius: "6px", cursor: "pointer", fontWeight: "bold", fontSize: "1rem", flex: 1, textDecoration: "none" };
const btnStyleAlert = { display: "flex", alignItems: "center", justifyContent: "center", gap: "0.5rem", padding: "0.75rem 1.25rem", backgroundColor: "#111827", color: "#FACC15", border: "none", borderRadius: "6px", cursor: "pointer", fontWeight: "bold", fontSize: "1rem" };